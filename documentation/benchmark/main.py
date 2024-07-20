import asyncio
import base64
import datetime
import fcntl
import hashlib
import random
import secrets
import socket
import struct
import time
from typing import Any, List, Literal, Mapping, Optional, TypedDict

import aiohttp
from aiohttp.client import ClientSession

TMethod = Literal["putflag", "getflag", "exploit", "putnoise", "getnoise", "havoc"]
TVariants = Mapping[int, List[TMethod]]


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

VARIANTS: TVariants = {
    0: ["putflag", "getflag", "exploit", "putnoise", "getnoise", "havoc"],
    1: ["putflag", "getflag", "exploit", "putnoise", "getnoise", "havoc"],
}
TICKS = 20
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
    methods: List[TMethod]

    def __init__(
        self,
        session: ClientSession,
        round_id: int,
        variant_id: int,
        methods: List[TMethod],
    ):
        self.session = session
        self.round_id = round_id
        self.variant_id = variant_id
        self.chain_prefix = secrets.token_hex(20)
        self.flag = generate_dummyflag()
        self.methods = methods

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
    def putnoise_chain_id(self) -> str:
        return f"{self.chain_prefix}_noise_s0_r{self.round_id}_t0_i0"

    @property
    def getnoise_chain_id(self) -> str:
        return f"{self.chain_prefix}_noise_s0_r{self.round_id}_t0_i0"

    @property
    def havoc_chain_id(self) -> str:
        return f"{self.chain_prefix}_havoc_s0_r{self.round_id}_t0_i0"

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

    @property
    def putnoise_payload(self) -> CheckerPayload:
        payload = self.payload_stub
        payload["method"] = "putnoise"
        payload["taskChainId"] = self.putnoise_chain_id
        return payload

    @property
    def getnoise_payload(self) -> CheckerPayload:
        payload = self.payload_stub
        payload["method"] = "getnoise"
        payload["taskChainId"] = self.getnoise_chain_id
        return payload

    @property
    def havoc_payload(self) -> CheckerPayload:
        payload = self.payload_stub
        payload["method"] = "havoc"
        payload["taskChainId"] = self.getnoise_chain_id
        return payload

    async def request(self, method: TMethod):
        if method not in self.methods:
            return None
        payload = None
        match method:
            case "putflag":
                payload = self.putflag_payload
            case "getflag":
                payload = self.getflag_payload
            case "exploit":
                payload = self.exploit_payload
            case "putnoise":
                payload = self.putnoise_payload
            case "getnoise":
                payload = self.getnoise_payload
            case "havoc":
                payload = self.havoc_payload
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
        variant_ids: TVariants,
        multiplier: int = 1,
        exploits_amount: int = 20,
    ):
        self.round_id = round_id
        self.variants = []
        for variant_id, methods in variant_ids.items():
            for _ in range(multiplier):
                self.variants.append(Variant(client, round_id, variant_id, methods))
        self.exploits_amount = exploits_amount

    async def request(self, method: TMethod):
        futures = []
        for variant in self.variants:
            futures.append(variant.request(method))

        return await asyncio.gather(*futures)

    async def exec(self, curr_round: int = -1, with_put: bool = False):
        if with_put:
            await self.request("putflag")
            await self.request("putnoise")
        coros = []
        if curr_round >= 0 and curr_round < self.round_id + 3:
            coros.append(run_in(30, self.request("getflag")))
        if curr_round >= 0 and curr_round < self.round_id + 3:
            coros.append(run_in(30, self.request("getnoise")))
        if curr_round < self.round_id + 3:
            coros.append(run_before(60, self.request("havoc")))
            coros.append(run_before(40, self.request("havoc")))
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
            await round.request("putnoise")
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

            tasks.append(asyncio.create_task(round.exec(with_put=True)))
            for _round in rounds[-EXPLOIT_PAST_ROUNDS:]:
                tasks.append(asyncio.create_task(_round.exec(curr_round)))
            rounds.append(round)

            await asyncio.sleep(60)
            print(f"TICK {i} END ({str(time.monotonic() - start)[:5]})")
            curr_round += 1

        await asyncio.gather(*tasks)


asyncio.run(main())
