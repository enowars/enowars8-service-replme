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
from typing import Any, List, Optional, TypedDict

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

VARIANTS = [0, 1]
TICKS = 3
MULTIPLIER = 1
EXPLOITS_AMOUNT = 20
EXPLOIT_PAST_ROUNDS = 10


async def run_before(delay: int, coro):
    _delay = random.randint(0, delay * 1000) / 1000
    await asyncio.sleep(_delay)
    await coro


async def run_in(delay: int, coro):
    await asyncio.sleep(delay)
    await coro


class CheckerPayload(TypedDict):
    address: str
    attackInfo: Optional[Any]
    currentRoundId: int
    flag: Optional[str]
    flagHash: Optional[str]
    flagRegex: Optional[str]
    method: Optional[str]
    relatedRoundId: int
    roundLength: int
    taskChainId: Optional[str]
    taskId: int
    teamId: int
    teamName: str
    timeout: int
    variantId: int


def generate_dummyflag() -> str:
    flag = "ENO" + base64.b64encode(secrets.token_bytes(36)).decode()
    assert len(flag) == 51
    return flag


class Variant:
    session: ClientSession
    round_id: int
    variant_id: int
    chain_prefix: str
    flag: str
    attack_info: str

    def __init__(self, session: ClientSession, round_id: int, variant_id: int):
        self.session = session
        self.round_id = round_id
        self.variant_id = variant_id
        self.chain_prefix = secrets.token_hex(20)
        self.flag = generate_dummyflag()

    @property
    def flag_hash(self) -> str:
        return hashlib.sha256(self.flag.encode()).hexdigest()

    @property
    def payload_stub(self) -> CheckerPayload:
        return {
            "address": SERVICE_ADDR,
            "attackInfo": None,
            "currentRoundId": self.round_id,
            "flag": None,
            "flagHash": None,
            "flagRegex": None,
            "method": None,
            "relatedRoundId": self.round_id,
            "roundLength": 60000,
            "taskChainId": None,
            "taskId": self.round_id,
            "teamId": 0,
            "teamName": "teamname",
            "timeout": 10000,
            "variantId": self.variant_id,
        }

    @property
    def putflag_chain_id(self) -> str:
        return f"{self.chain_prefix}_flag_s0_r{self.round_id}_t0_i0"

    @property
    def getflag_chain_id(self) -> str:
        return f"{self.chain_prefix}_flag_s0_r{self.round_id}_t0_i0"

    @property
    def exploit_chain_id(self) -> str:
        return f"{self.chain_prefix}_exploit_s0_r{self.round_id}_t0_i0"

    @property
    def putflag_payload(self) -> CheckerPayload:
        payload = self.payload_stub
        payload["flag"] = self.flag
        payload["method"] = "putflag"
        payload["taskChainId"] = self.putflag_chain_id
        return payload

    @property
    def getflag_payload(self) -> CheckerPayload:
        payload = self.payload_stub
        payload["flag"] = self.flag
        payload["method"] = "getflag"
        payload["taskChainId"] = self.getflag_chain_id
        return payload

    @property
    def exploit_payload(self) -> CheckerPayload:
        payload = self.payload_stub
        payload["attackInfo"] = self.attack_info
        payload["flagHash"] = self.flag_hash
        payload["flagRegex"] = "ENO[A-Za-z0-9+\\/=]{48}"
        payload["method"] = "exploit"
        payload["taskChainId"] = self.exploit_chain_id
        return payload

    async def request(self, method: str):
        payload = None
        match method:
            case "putflag":
                payload = self.putflag_payload
            case "getflag":
                payload = self.getflag_payload
            case "exploit":
                payload = self.exploit_payload
        start = time.monotonic()
        response = await self.session.post("/", json=payload)
        if response.status == 200:
            payload = await response.json()
            s = f"{datetime.datetime.now()} {method.upper()} {self.round_id}:{self.variant_id}: ({str(time.monotonic() - start)[:5]}) {payload['result']} - {payload['message']}"
            print(s)
            if method == "putflag":
                self.attack_info = payload["attackInfo"]
        else:
            s = f"{datetime.datetime.now()} {method.upper()} {self.round_id}:{self.variant_id}: ({str(time.monotonic() - start)[:5]}) {response.status}"
        return response


class Round:
    round_id: int
    variants: List[Variant]
    exploits_amount: int

    def __init__(
        self,
        client: ClientSession,
        round_id: int,
        variant_ids: List[int],
        multiplier: int = 1,
        exploits_amount: int = 20,
    ):
        self.round_id = round_id
        self.variants = []
        for variant_id in variant_ids:
            for _ in range(multiplier):
                self.variants.append(Variant(client, round_id, variant_id))
        self.exploits_amount = exploits_amount

    async def request(self, method: str):
        futures = []
        for variant in self.variants:
            futures.append(variant.request(method))

        return await asyncio.gather(*futures)

    async def exec(self, curr_round: int = -1, with_putflag: bool = False):
        if with_putflag:
            await self.request("putflag")
        coros = []
        if curr_round >= 0 and curr_round < self.round_id + 3:
            coros.append(run_in(30, self.request("getflag")))
        for variant in self.variants:
            for _ in range(self.exploits_amount):
                coros.append(run_before(60, variant.request("exploit")))
        return await asyncio.gather(*coros)


async def main():
    curr_round = 0
    rounds: List[Round] = []

    tasks: List[asyncio.Task] = []

    async with aiohttp.ClientSession(CHECKER_ADDR) as client:
        for i in range(10):
            round = Round(
                client,
                curr_round,
                VARIANTS,
                multiplier=MULTIPLIER,
                exploits_amount=EXPLOITS_AMOUNT,
            )
            await round.request("putflag")
            rounds.append(round)
            curr_round += 1

        for i in range(TICKS):
            print("TICK", i)
            start = time.monotonic()
            round = Round(
                client,
                curr_round,
                VARIANTS,
                multiplier=MULTIPLIER,
                exploits_amount=EXPLOITS_AMOUNT,
            )

            tasks.append(asyncio.create_task(round.exec(with_putflag=True)))
            for _round in rounds[-EXPLOIT_PAST_ROUNDS:]:
                tasks.append(asyncio.create_task(_round.exec(curr_round)))

            await asyncio.sleep(60)
            print(f"TICK {i} END ({str(time.monotonic() - start)[:5]})")
            curr_round += 1

        await asyncio.gather(*tasks)


asyncio.run(main())
