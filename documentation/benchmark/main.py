import asyncio
import base64
import fcntl
import hashlib
import secrets
import socket
import struct
import time
import random
import datetime
from typing import Any, List, Mapping

import aiohttp
from aiohttp.client import ClientSession


def get_ip_address(ifname):
    with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as s:
        return socket.inet_ntoa(
            fcntl.ioctl(
                s.fileno(),
                0x8915,  # SIOCGIFADDR
                struct.pack("256s", bytes(ifname[:15], "utf-8")),
            )[20:24]
        )


CHECKER_ADDR = "http://127.0.0.1:16969"
SERVICE_ADDR = get_ip_address("wlp3s0")

"""
200
{
  "result": "OK",
  "message": "string",
  "attackInfo": "string",
  "flag": "string"
}
"""

"""
422
{
  "detail": [
    {
      "loc": [
        "string",
        0
      ],
      "msg": "string",
      "type": "string"
    }
  ]
}
"""


def generate_dummyflag() -> str:
    flag = "ENO" + base64.b64encode(secrets.token_bytes(36)).decode()
    assert len(flag) == 51
    return flag


class Round:
    round_id: int
    chain_prefix: str
    flag0: str
    flag1: str
    attack_info0: str
    attack_info1: str

    def __init__(self, round_id: int):
        self.round_id = round_id
        self.chain_prefix = secrets.token_hex(20)
        self.flag0 = generate_dummyflag()
        self.flag1 = generate_dummyflag()
        self.attack_info0 = ""
        self.attack_info1 = ""

    @property
    def putflag0_chain_id(self) -> str:
        return f"{self.chain_prefix}_flag_s0_r{self.round_id}_t0_i0"

    @property
    def getflag0_chain_id(self) -> str:
        return f"{self.chain_prefix}_flag_s0_r{self.round_id}_t0_i0"

    @property
    def exploit0_chain_id(self) -> str:
        return f"{self.chain_prefix}_exploit_s0_r{self.round_id}_t0_i0"

    @property
    def putflag1_chain_id(self) -> str:
        return f"{self.chain_prefix}_flag_s0_r{self.round_id}_t0_i1"

    @property
    def getflag1_chain_id(self) -> str:
        return f"{self.chain_prefix}_flag_s0_r{self.round_id}_t0_i1"

    @property
    def exploit1_chain_id(self) -> str:
        return f"{self.chain_prefix}_exploit_s0_r{self.round_id}_t0_i1"

    @property
    def putflag_payload0(self) -> Mapping[str, Any]:
        return {
            "address": SERVICE_ADDR,
            "attackInfo": None,
            "currentRoundId": self.round_id,
            "flag": self.flag0,
            "flagHash": None,
            "flagRegex": None,
            "method": "putflag",
            "relatedRoundId": self.round_id,
            "roundLength": 60000,
            "taskChainId": self.putflag0_chain_id,
            "taskId": self.round_id,
            "teamId": 0,
            "teamName": "teamname",
            "timeout": 10000,
            "variantId": 0,
        }

    @property
    def getflag_payload0(self) -> Mapping[str, Any]:
        return {
            "address": SERVICE_ADDR,
            "attackInfo": None,
            "currentRoundId": self.round_id,
            "flag": self.flag0,
            "flagHash": None,
            "flagRegex": None,
            "method": "getflag",
            "relatedRoundId": self.round_id,
            "roundLength": 60000,
            "taskChainId": self.getflag0_chain_id,
            "taskId": self.round_id,
            "teamId": 0,
            "teamName": "teamname",
            "timeout": 10000,
            "variantId": 0,
        }

    @property
    def exploit_payload0(self) -> Mapping[str, Any]:
        return {
            "address": SERVICE_ADDR,
            "attackInfo": self.attack_info0,
            "currentRoundId": self.round_id,
            "flag": None,
            "flagHash": hashlib.sha256(self.flag0.encode()).hexdigest(),
            "flagRegex": "ENO[A-Za-z0-9+\\/=]{48}",
            "method": "exploit",
            "relatedRoundId": self.round_id,
            "roundLength": 60000,
            "taskChainId": self.exploit0_chain_id,
            "taskId": self.round_id,
            "teamId": 0,
            "teamName": "teamname",
            "timeout": 10000,
            "variantId": 0,
        }

    @property
    def putflag_payload1(self) -> Mapping[str, Any]:
        return {
            "address": SERVICE_ADDR,
            "attackInfo": None,
            "currentRoundId": self.round_id,
            "flag": self.flag0,
            "flagHash": None,
            "flagRegex": None,
            "method": "putflag",
            "relatedRoundId": self.round_id,
            "roundLength": 60000,
            "taskChainId": self.putflag1_chain_id,
            "taskId": self.round_id,
            "teamId": 0,
            "teamName": "teamname",
            "timeout": 10000,
            "variantId": 1,
        }

    @property
    def getflag_payload1(self) -> Mapping[str, Any]:
        return {
            "address": SERVICE_ADDR,
            "attackInfo": None,
            "currentRoundId": self.round_id,
            "flag": self.flag0,
            "flagHash": None,
            "flagRegex": None,
            "method": "getflag",
            "relatedRoundId": self.round_id,
            "roundLength": 60000,
            "taskChainId": self.getflag1_chain_id,
            "taskId": self.round_id,
            "teamId": 0,
            "teamName": "teamname",
            "timeout": 10000,
            "variantId": 1,
        }

    @property
    def exploit_payload1(self) -> Mapping[str, Any]:
        return {
            "address": SERVICE_ADDR,
            "attackInfo": self.attack_info1,
            "currentRoundId": self.round_id,
            "flag": None,
            "flagHash": hashlib.sha256(self.flag0.encode()).hexdigest(),
            "flagRegex": "ENO[A-Za-z0-9+\\/=]{48}",
            "method": "exploit",
            "relatedRoundId": self.round_id,
            "roundLength": 60000,
            "taskChainId": self.exploit1_chain_id,
            "taskId": self.round_id,
            "teamId": 0,
            "teamName": "teamname",
            "timeout": 10000,
            "variantId": 1,
        }


errs = []


async def putflag0(client: ClientSession, round: Round):
    global errs
    start = time.monotonic()
    async with client.post(
        CHECKER_ADDR,
        json=round.putflag_payload0,
    ) as response:
        if response.status == 200:
            payload = await response.json()
            s = f"{datetime.datetime.now()} PUTFLAG {round.round_id}:0: ({str(time.monotonic() - start)[:5]}) {payload['result']} - {payload['message']}"
            print(s)
            if payload["result"] != "OK":
                errs.append(s)
            round.attack_info0 = payload["attackInfo"]
        else:
            s = f"{datetime.datetime.now()} PUTFLAG {round.round_id}:0: ({str(time.monotonic() - start)[:5]}) {response.status}"
            print(s)
            errs.append(s)


async def getflag0(client: ClientSession, round: Round):
    global errs
    start = time.monotonic()
    async with client.post(
        CHECKER_ADDR,
        json=round.getflag_payload0,
    ) as response:
        if response.status == 200:
            payload = await response.json()
            s = f"{datetime.datetime.now()} GETFLAG {round.round_id}:0: ({str(time.monotonic() - start)[:5]}) {payload['result']} - {payload['message']}"
            print(s)
            if payload["result"] != "OK":
                errs.append(s)
        else:
            s = f"{datetime.datetime.now()} GETFLAG {round.round_id}:0: ({str(time.monotonic() - start)[:5]}) {response.status}"
            print(s)
            errs.append(s)


async def exploit0(client: ClientSession, round: Round):
    global errs
    start = time.monotonic()
    async with client.post(
        CHECKER_ADDR,
        json=round.exploit_payload0,
    ) as response:
        if response.status == 200:
            payload = await response.json()
            s = f"{datetime.datetime.now()} EXPLOIT {round.round_id}:0: ({str(time.monotonic() - start)[:5]}) {payload['result']} - {payload['message']}"
            print(s)
            if payload["result"] != "OK":
                errs.append(s)
        else:
            s = f"{datetime.datetime.now()} EXPLOIT {round.round_id}:0: ({str(time.monotonic() - start)[:5]}) {response.status}"
            print(s)
            errs.append(s)


async def putflag1(client: ClientSession, round: Round):
    global errs
    start = time.monotonic()
    async with client.post(
        CHECKER_ADDR,
        json=round.putflag_payload1,
    ) as response:
        if response.status == 200:
            payload = await response.json()
            s = f"{datetime.datetime.now()} PUTFLAG {round.round_id}:1: ({str(time.monotonic() - start)[:5]}) {payload['result']} - {payload['message']}"
            print(s)
            if payload["result"] != "OK":
                errs.append(s)
            round.attack_info1 = payload["attackInfo"]
        else:
            s = f"{datetime.datetime.now()} PUTFLAG {round.round_id}:1: ({str(time.monotonic() - start)[:5]}) {response.status}"
            print(s)
            errs.append(s)


async def getflag1(client: ClientSession, round: Round):
    global errs
    start = time.monotonic()
    async with client.post(
        CHECKER_ADDR,
        json=round.getflag_payload1,
    ) as response:
        if response.status == 200:
            payload = await response.json()
            s = f"{datetime.datetime.now()} GETFLAG {round.round_id}:1: ({str(time.monotonic() - start)[:5]}) {payload['result']} - {payload['message']}"
            print(s)
            if payload["result"] != "OK":
                errs.append(s)
        else:
            s = f"{datetime.datetime.now()} GETFLAG {round.round_id}:1: ({str(time.monotonic() - start)[:5]}) {response.status}"
            print(s)
            errs.append(s)


async def exploit1(client: ClientSession, round: Round):
    global errs
    start = time.monotonic()
    async with client.post(
        CHECKER_ADDR,
        json=round.exploit_payload1,
    ) as response:
        if response.status == 200:
            payload = await response.json()
            s = f"{datetime.datetime.now()} EXPLOIT {round.round_id}:1: ({str(time.monotonic() - start)[:5]}) {payload['result']} - {payload['message']}"
            print(s)
            if payload["result"] != "OK":
                errs.append(s)
        else:
            s = f"{datetime.datetime.now()} EXPLOIT {round.round_id}:1: ({str(time.monotonic() - start)[:5]}) {response.status}"
            print(s)
            errs.append(s)


async def exploit0_delayed(client: ClientSession, round: Round):
    delay = random.randint(0, 5500) / 100
    await asyncio.sleep(delay)
    await exploit0(client, round)


async def exploit1_delayed(client: ClientSession, round: Round):
    delay = random.randint(0, 5500) / 100
    await asyncio.sleep(delay)
    await exploit1(client, round)


async def exec_round(client: ClientSession, round: Round, with_putflag: bool = False):
    if with_putflag:
        await putflag0(client, round)
        await putflag1(client, round)
    getflag0_future = getflag0(client, round)
    getflag1_future = getflag1(client, round)
    futures0 = [exploit0_delayed(client, round) for _ in range(20)]
    futures1 = [exploit1_delayed(client, round) for _ in range(40)]
    await asyncio.gather(getflag0_future, getflag1_future, *futures0, *futures1)


async def main():
    global errs
    currentRound = 0
    rounds: List[Round] = []

    async with aiohttp.ClientSession() as client:
        for i in range(10):
            round = Round(currentRound)
            await putflag0(client, round)
            await putflag1(client, round)
            rounds.append(round)
            currentRound += 1

        for i in range(10):
            print("TICK", i)
            start = time.monotonic()
            round = Round(currentRound)
            current_round_future = exec_round(client, round, True)
            prev_rounds_futures = [
                exec_round(client, prevRound) for prevRound in rounds[-10:]
            ]
            rounds.append(round)
            await asyncio.gather(current_round_future, *prev_rounds_futures)
            print(f"TICK {i} END ({str(time.monotonic() - start)[:5]})")
            currentRound += 1
    print("============ Result ============")
    for err in errs:
        print(err)


asyncio.run(main())
