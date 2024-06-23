import string
import random
import asyncio
import re
from logging import LoggerAdapter
from typing import List, Optional
import websockets

from enochecker3.chaindb import ChainDB
from enochecker3.types import MumbleException
from enochecker3.utils import assert_equals
from httpx import AsyncClient, Cookies


class ShellCommand:
    _cmd: str
    ok: str
    err: str

    def __init__(self, _cmd: str, ok: str, err: str) -> None:
        self.cmd = _cmd
        self.ok = ok
        self.err = err

    @property
    def cmd(self) -> str:
        return self._cmd + "\n"

    @cmd.setter
    def cmd(self, _cmd):
        self._cmd = _cmd


class ShellCommandErrorBuilder:
    cmd: str
    ok: str

    def __init__(self, cmd: str, ok: str) -> None:
        self.cmd = cmd
        self.ok = ok

    def err(self) -> ShellCommand:
        return ShellCommand(self.cmd + ' || echo "\\nERROR"', self.ok, ".*\nERROR.*")

    def errext(self) -> ShellCommand:
        return ShellCommand(
            self.cmd + ' && echo "\\nERROR" || echo "\\nERROR"', self.ok, ".*\nERROR.*"
        )

    def fail(self, err: str, suffix="") -> ShellCommand:
        return ShellCommand(self.cmd + suffix, self.ok, err)

    def default(self) -> ShellCommand:
        return self.err()


class ShellCommandBuilder:
    cmd: str

    def __init__(self, cmd: str) -> None:
        self.cmd = cmd

    def ok(self) -> ShellCommandErrorBuilder:
        return ShellCommandErrorBuilder(self.cmd + ' && echo "\\nOK"', ".*\nOK.*")

    def expect(self, ok: str, suffix="") -> ShellCommandErrorBuilder:
        return ShellCommandErrorBuilder(self.cmd + suffix, ok)

    def default(self) -> ShellCommand:
        return self.ok().err()


def sh(cmd: str):
    return ShellCommandBuilder(cmd)


class ShellCommandChain:
    command_chain: List[ShellCommand]
    validation_chain: List[ShellCommand]

    def __init__(
        self, cmds: List[ShellCommand], validations: List[ShellCommand]
    ) -> None:
        self.command_chain = cmds
        self.validation_chain = validations


def shchain(cmds: List[ShellCommand] = [], validations: List[ShellCommand] = []):
    return ShellCommandChain(cmds, validations)


async def do_user_auth(
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
    id = json["id"]
    assert_equals(id is not None, True, "Did not receive a repl id")
    logger.info(f"REPL-ID: {id}")

    cookies = response.cookies

    for k, v in cookies.items():
        logger.info(f"Cookie: {k}={v}")

    return (cookies, id)


async def user_create(
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
) -> tuple[str, Cookies, str]:
    username = "".join(random.choice(string.ascii_lowercase) for _ in range(60))
    password = "".join(random.choice(string.ascii_lowercase) for _ in range(60))

    logger.info(f"Creating user: {username}:{password}")
    (cookies, id) = await do_user_auth(client, logger, username, password)
    await db.set("credentials", (username, password))

    return (username, cookies, id)


async def user_login(
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
) -> tuple[Cookies, str]:
    try:
        (username, password) = await db.get("credentials")
    except KeyError:
        raise MumbleException("Missing database entry from putflag")

    logger.info(f"Authenticating user: {username}:{password}")
    (cookies, id) = await do_user_auth(client, logger, username, password)

    return (cookies, id)


def is_list_of_str(element):
    if isinstance(element, list):
        return all(isinstance(item, str) for item in element)
    return False


async def get_sessions(
    client: AsyncClient,
    cookies: Cookies,
    logger: LoggerAdapter,
) -> List[str]:
    response = await client.get(
        "/api/user/sessions",
        follow_redirects=True,
        headers={"Cookie": "session=" + (cookies.get("session") or "")},
    )
    logger.info(f"Server answered: {response.status_code} - {response.text}")
    assert_equals(response.status_code < 300, True, "Getting sessions failed")

    json = response.json()
    assert_equals(is_list_of_str(json), True, "Did not receive valid sessions obj")
    logger.info(f"Sessions: {json}")

    return json


async def websocket_recv_until(
    websocket: websockets.WebSocketClientProtocol,
    expected: str,
    unexpected: Optional[str],
    logger: Optional[LoggerAdapter] = None,
) -> str:
    payload = ""
    match = None

    while match is None:
        msg = await asyncio.wait_for(websocket.recv(), timeout=10)
        if isinstance(msg, bytes):
            payload += msg.decode("utf-8")
        elif isinstance(msg, str):
            payload += msg
        else:
            raise MumbleException("Websocket connection recv returned unexpected data")
        if logger is not None:
            logger.info("PAYLOAD: " + payload)
        match = re.match(expected, payload, re.S)
        if (
            match is None
            and unexpected is not None
            and re.match(unexpected, payload, re.S) is not None
        ):
            raise MumbleException("Received unexpected input")

    return payload


async def terminal_websocket(
    address: str,
    logger: LoggerAdapter,
    cookies: Cookies,
    id: str,
    actions: List[ShellCommand],
):
    url = f"ws://{address}:6969/api/repl/{id}"
    response = ""
    cookie = cookies.get("session")
    async with websockets.connect(
        url, extra_headers={"Cookie": f"session={cookie}"}
    ) as websocket:
        await websocket_recv_until(websocket, ".*%.*", None)
        for action in actions:
            logger.info("Sending: " + str(action.cmd.encode("utf-8")))
            await websocket.send(action.cmd)
            logger.info("Expecting: " + str(action.ok.encode("utf-8")))
            _response = await websocket_recv_until(
                websocket, action.ok, action.err, logger
            )
            logger.info("Got: " + _response)
            response += _response

    return response
