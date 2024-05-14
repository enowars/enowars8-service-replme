from asyncio import sleep
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

from util import create_terminal, create_user, terminal_websocket, user_login

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
    (username, password, port) = await create_user(client, db, logger)
    pid = await create_terminal(client, logger, username, password)
    await terminal_websocket(
        task.address,
        logger,
        username,
        password,
        port,
        pid,
        [
            (f'echo "{task.flag}" > flagstore.txt && echo OK\n', ".*\nOK.*"),
        ],
    )

    return username


@checker.getflag(0)
async def getflag_test(
    task: GetflagCheckerTaskMessage,
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
) -> None:
    (username, password, port) = await user_login(client, db, logger)
    pid = await create_terminal(client, logger, username, password)
    await terminal_websocket(
        task.address,
        logger,
        username,
        password,
        port,
        pid,
        [
            ("cat flagstore.txt\n", f".*\n{task.flag}.*"),
        ],
    )


"""
@checker.exploit(0)
async def exploit_test(searcher: FlagSearcher, client: AsyncClient) -> Optional[str]:
    r = await client.get(
        "/note/*",
    )
    assert not r.is_error

    if flag := searcher.search_flag(r.text):
        return flag
"""


if __name__ == "__main__":
    checker.run()
