import string
import random
import asyncio
import base64
import re
from logging import LoggerAdapter
from typing import List, Optional
import websockets

from enochecker3.chaindb import ChainDB
from enochecker3.types import MumbleException
from enochecker3.utils import assert_equals
from httpx import AsyncClient, Cookies, Response


async def create_user(
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
) -> tuple[str, str, Cookies, str]:
    username = "".join(random.choice(string.ascii_lowercase) for _ in range(60))
    password = "".join(random.choice(string.ascii_lowercase) for _ in range(60))

    logger.info(f"Creating user: {username}:{password}")
    response = await client.post(
        "/api/repl",
        data={"username": username, "password": password},
        follow_redirects=True,
    )
    logger.info(f"Server answered: {response.status_code} - {response.text}")

    assert_equals(response.status_code < 300, True, "Creating user failed")

    json = response.json()
    replUuid = json["replUuid"]
    assert_equals(replUuid is not None, True, "Did not receive a replUuid")
    logger.info(f"Port: {replUuid}")

    await db.set("credentials", (username, password))

    return (username, password, response.cookies, replUuid)


async def do_user_login(
    client: AsyncClient, logger: LoggerAdapter, username: str, password: str
) -> tuple[Cookies, str]:
    response = await client.post(
        "/api/repl",
        data={"username": username, "password": password},
        follow_redirects=True,
    )
    logger.info(f"Server answered: {response.status_code} - {response.text}")

    assert_equals(response.status_code < 300, True, "Creating user failed")

    json = response.json()
    replUuid = json["replUuid"]
    assert_equals(replUuid is not None, True, "Did not receive a replUuid")
    logger.info(f"Port: {replUuid}")

    return (response.cookies, replUuid)


async def user_login(
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
) -> tuple[str, str, Cookies, str]:
    try:
        (username, password) = await db.get("credentials")
    except KeyError:
        raise MumbleException("Missing database entry from putflag")

    (cookies, replUuid) = await do_user_login(client, logger, username, password)

    return (username, password, cookies, replUuid)


async def create_terminal(
    client: AsyncClient, logger: LoggerAdapter, username: str, password: str
) -> str:
    logger.info("Creating terminal")
    response: Optional[Response] = None
    for _ in range(5):
        response = await client.post(
            "/api/terminal",
            headers={
                "Authorization": "Basic "
                + base64.b64encode(f"{username}:{password}".encode("utf-8")).decode(
                    "utf-8"
                )
            },
            json={"cols": 80, "rows": 25},
            follow_redirects=True,
        )
        if response.status_code < 300:
            break

    if response is None:
        raise MumbleException("Create terminal request failed")

    logger.info(f"Server answered: {response.status_code} - {response.text}")

    assert_equals(response.status_code < 300, True, "Creating terminal failed")

    pid = response.text
    assert_equals(pid is not None, True, "Did not receive a PID")
    logger.info(f"PID: {pid}")

    return pid


async def websocket_recv_until(
    websocket: websockets.WebSocketClientProtocol,
    pattern: str,
    logger: Optional[LoggerAdapter] = None,
) -> str:
    payload = ""
    match = None

    while match is None:
        msg = await asyncio.wait_for(websocket.recv(), timeout=0.5)
        if isinstance(msg, bytes):
            payload += msg.decode("utf-8")
        elif isinstance(msg, str):
            payload += msg
        else:
            raise MumbleException("Websocket connection recv returned unexpected data")
        if logger is not None:
            logger.info("PAYLOAD: " + payload)
        match = re.match(pattern, payload, re.S)

    return payload


async def terminal_login(
    logger: LoggerAdapter,
    websocket: websockets.WebSocketClientProtocol,
    username: str,
    password: str,
):
    logger.info(await websocket_recv_until(websocket, ".*login:.*"))
    await websocket.send(username + "\n")
    logger.info(await websocket_recv_until(websocket, ".*Password:.*"))
    await websocket.send(password + "\n")
    logger.info(await websocket_recv_until(websocket, ".*Welcome to Alpine!.*"))


async def terminal_websocket(
    address: str,
    logger: LoggerAdapter,
    cookies: Cookies,
    replUuid: str,
    actions: List[tuple[str, str]],
):
    url = f"ws://{address}:6969/api/repl/{replUuid}"
    response = ""
    async with websockets.connect(
        url, extra_headers={"Cookie": f"session={cookies.get("session")}"}
    ) as websocket:
        await websocket_recv_until(websocket, ".*%.*")
        for action in actions:
            logger.info("Sending: " + str(action[0].encode("utf-8")))
            await websocket.send(action[0])
            logger.info("Expecting: " + str(action[1].encode("utf-8")))
            _response = await websocket_recv_until(websocket, action[1], logger)
            logger.info("Got: " + _response)
            response += _response

    return response
