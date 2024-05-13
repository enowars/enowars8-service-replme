from asyncio import sleep
import secrets
import base64
import websockets
from logging import LoggerAdapter
from typing import Optional

from enochecker3.chaindb import ChainDB
from enochecker3.enochecker import Enochecker
from enochecker3.types import (
    GetflagCheckerTaskMessage,
    MumbleException,
    PutflagCheckerTaskMessage,
)
from enochecker3.utils import FlagSearcher, assert_equals, assert_in
from httpx import AsyncClient

checker = Enochecker("cafedodo", 6969)


def app():
    return checker.app


@checker.putflag(0)
async def putflag_test(
    task: PutflagCheckerTaskMessage,
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
) -> str:
    username = secrets.token_hex(30)
    password = secrets.token_hex(30)

    await db.set("credentials", (username, password))

    logger.info(f"Creating user: {username}:{password}")
    response = await client.post(
        "/api/login/private",
        data={"username": username, "password": password},
        follow_redirects=True,
    )
    logger.info(f"Server answered: {response.status_code} - {response.text}")

    assert_equals(response.status_code < 300, True, "Creating user failed")

    json = response.json()
    port = json["port"]
    assert_equals(port is not None, True, "Did not receive a port")
    logger.info(f"Port: {port}")

    logger.info("Creating terminal")
    for _ in range(5):
        response = await client.post(
            "/api/terminal",
            headers={
                "Authorization": "Basic "
                + base64.b64encode(f"{username}:{password}".encode("utf-8")).decode("utf-8")
            },
            json={"cols": 80, "rows": 25},
            follow_redirects=True,
        )
        if response.status_code < 300:
            break
    logger.info(f"Server answered: {response.status_code} - {response.text}")

    assert_equals(response.status_code < 300, True, "Creating terminal failed")

    pid = response.text
    assert_equals(pid is not None, True, "Did not receive a PID")
    logger.info(f"PID: {pid}")

    url = f"ws://{task.address}:6969/ws/terminal/{port}/{pid}"
    async with websockets.connect(url) as websocket:
        await sleep(0.1)
        await websocket.send(username + '\n')
        await sleep(0.1)
        await websocket.send(password + '\n')
        await sleep(0.1)
        message = await websocket.recv()
        while message is not None:
            logger.info(f"Message: {message}")
            message = await websocket.recv()

    return username


@checker.getflag(0)
async def getflag_test(
    task: GetflagCheckerTaskMessage, client: AsyncClient, db: ChainDB
) -> None:
    try:
        token = await db.get("token")
    except KeyError:
        raise MumbleException("Missing database entry from putflag")

    r = await client.get(f"/note/{token}")
    assert_equals(r.status_code, 200, "getting note with flag failed")
    assert_in(task.flag, r.text, "flag missing from note")


@checker.exploit(0)
async def exploit_test(searcher: FlagSearcher, client: AsyncClient) -> Optional[str]:
    r = await client.get(
        "/note/*",
    )
    assert not r.is_error

    if flag := searcher.search_flag(r.text):
        return flag


if __name__ == "__main__":
    checker.run()
