import asyncio
import base64
import fcntl
import hashlib
import secrets
import socket
import struct
import time
from typing import Any, Mapping

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
    flag: str
    attack_info: str

    def __init__(self, round_id: int):
        self.round_id = round_id
        self.chain_prefix = secrets.token_hex(20)
        self.flag = generate_dummyflag()
        self.attack_info = ""

    @property
    def putflag_chain_id(self) -> str:
        return f"{self.chain_prefix}_flag_s0_r{self.round_id}_t0_i0"

    @property
    def getflag_chain_id(self) -> str:
        return f"{self.chain_prefix}_flag_s0_r{self.round_id}_t0_i0"

    @property
    def exploit_chain_id(self) -> str:
        return f"{self.chain_prefix}_exploit{self.round_id}_t0_i0"

    @property
    def putflag_payload(self) -> Mapping[str, Any]:
        return {
            "address": SERVICE_ADDR,
            "attackInfo": None,
            "currentRoundId": self.round_id,
            "flag": self.flag,
            "flagHash": None,
            "flagRegex": None,
            "method": "putflag",
            "relatedRoundId": self.round_id,
            "roundLength": 60000,
            "taskChainId": self.putflag_chain_id,
            "taskId": self.round_id,
            "teamId": 0,
            "teamName": "teamname",
            "timeout": 10000,
            "variantId": 0,
        }

    @property
    def getflag_payload(self) -> Mapping[str, Any]:
        return {
            "address": SERVICE_ADDR,
            "attackInfo": None,
            "currentRoundId": self.round_id,
            "flag": self.flag,
            "flagHash": None,
            "flagRegex": None,
            "method": "getflag",
            "relatedRoundId": self.round_id,
            "roundLength": 60000,
            "taskChainId": self.getflag_chain_id,
            "taskId": self.round_id,
            "teamId": 0,
            "teamName": "teamname",
            "timeout": 10000,
            "variantId": 0,
        }

    @property
    def exploit_payload(self) -> Mapping[str, Any]:
        return {
            "address": SERVICE_ADDR,
            "attackInfo": self.attack_info,
            "currentRoundId": self.round_id,
            "flag": None,
            "flagHash": hashlib.sha256(self.flag.encode()).hexdigest(),
            "flagRegex": "ENO[A-Za-z0-9+\\/=]{48}",
            "method": "exploit",
            "relatedRoundId": self.round_id,
            "roundLength": 60000,
            "taskChainId": self.exploit_chain_id,
            "taskId": self.round_id,
            "teamId": 0,
            "teamName": "teamname",
            "timeout": 10000,
            "variantId": 0,
        }


async def putflag(client: ClientSession, round: Round):
    start = time.monotonic()
    async with client.post(
        CHECKER_ADDR,
        json=round.putflag_payload,
    ) as response:
        if response.status == 200:
            payload = await response.json()
            print(
                f"PUTFLAG {round.round_id}: ({str(time.monotonic() - start)[:5]}) {payload['result']} - {payload['message']}"
            )
            round.attack_info = payload["attackInfo"]
        else:
            print(
                f"PUTFLAG {round.round_id}: ({str(time.monotonic() - start)[:5]}) {response.status}"
            )


async def getflag(client: ClientSession, round: Round):
    start = time.monotonic()
    async with client.post(
        CHECKER_ADDR,
        json=round.getflag_payload,
    ) as response:
        if response.status == 200:
            payload = await response.json()
            print(
                f"GETFLAG {round.round_id}: ({str(time.monotonic() - start)[:5]}) {payload['result']} - {payload['message']}"
            )
        else:
            print(
                f"GETFLAG {round.round_id}: ({str(time.monotonic() - start)[:5]}) {response.status}"
            )


async def exploit(client: ClientSession, round: Round):
    start = time.monotonic()
    async with client.post(
        CHECKER_ADDR,
        json=round.exploit_payload,
    ) as response:
        if response.status == 200:
            payload = await response.json()
            print(
                f"EXPLOIT {round.round_id}: ({str(time.monotonic() - start)[:5]}) {payload['result']} - {payload['message']}"
            )
        else:
            print(
                f"EXPLOIT {round.round_id}: ({str(time.monotonic() - start)[:5]}) {response.status}"
            )


async def main():
    async with aiohttp.ClientSession() as client:
        round = Round(0)
        await putflag(client, round)
        await getflag(client, round)
        futures = [exploit(client, round) for _ in range(100)]
        await asyncio.gather(*futures)


asyncio.run(main())
