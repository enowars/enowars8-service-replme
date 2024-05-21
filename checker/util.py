import asyncio
import base64
import secrets
import re
from logging import LoggerAdapter
from typing import List, Optional
import websockets

from enochecker3.chaindb import ChainDB
from enochecker3.types import MumbleException
from enochecker3.utils import assert_equals
from httpx import AsyncClient, Response


async def create_user(
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
) -> tuple[str, str, str]:
    username = secrets.token_hex(30)
    password = secrets.token_hex(30)

    logger.info(f"Creating user: {username}:{password}")
    response = await client.post(
        "/api/login/private",
        data={"username": username, "password": password},
        follow_redirects=True,
    )
    logger.info(f"Server answered: {response.status_code} - {response.text}")

    assert_equals(response.status_code < 300, True, "Creating user failed")

    json = response.json()
    port = str(json["port"])
    assert_equals(port is not None, True, "Did not receive a port")
    logger.info(f"Port: {port}")

    await db.set("credentials", (username, password))

    return (username, password, port)


async def do_user_login(
    client: AsyncClient, logger: LoggerAdapter, username: str, password: str
) -> str:
    logger.info(f"Login user: {username}:{password}")
    response = await client.post(
        "/api/login/private",
        data={"username": username, "password": password},
        follow_redirects=True,
    )
    logger.info(f"Server answered: {response.status_code} - {response.text}")

    assert_equals(response.status_code < 300, True, "Creating user failed")

    json = response.json()
    port = str(json["port"])
    assert_equals(port is not None, True, "Did not receive a port")
    logger.info(f"Port: {port}")

    return port


async def user_login(
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
) -> tuple[str, str, str]:
    try:
        (username, password) = await db.get("credentials")
    except KeyError:
        raise MumbleException("Missing database entry from putflag")

    return (username, password, await do_user_login(client, logger, username, password))


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
    username: str,
    password: str,
    port: str,
    pid: str,
    actions: List[tuple[str, str]],
):
    url = f"ws://{address}:6969/ws/terminal/{port}/{pid}"
    response = ""
    async with websockets.connect(url) as websocket:
        await terminal_login(logger, websocket, username, password)
        for action in actions:
            logger.info("Sending: " + str(action[0].encode("utf-8")))
            await websocket.send(action[0])
            logger.info("Expecting: " + str(action[1].encode("utf-8")))
            _response = await websocket_recv_until(websocket, action[1], logger)
            logger.info("Got: " + _response)
            response += _response

    return response
