import logging
import time
from typing import List

import requests
from dataclasses import dataclass


@dataclass
class Vehicle:
    id: int
    name: str
    nation: str
    type: str
    level: str
    img: str
    current_count: int
    price: int
    next_price: int = None

    def __str__(self):
        return f'''{self.name}, {self.level} {self.nation} {self.type}
count left: {self.current_count}
current price: {self.price}; next price: {self.next_price if self.next_price else "-"}'''


class Poller:
    URL = 'https://eu.wotblitz.com/en/api/events/items/auction/?page_size=80&type[]=vehicle&saleable=true'

    def __init__(self, cache_timeout: int = 300):
        self._cache_timeout = cache_timeout
        self._last_polled_data: List[Vehicle] = []
        self._last_polled_at: int = 0
        self.poll_data()

    @staticmethod
    def parse_vehicle_data(vehicle) -> Vehicle:
        return Vehicle(id=vehicle['id'],
                       name=vehicle['entity']['user_string'],
                       nation=vehicle['entity']['nation'],
                       type=vehicle['entity']['type_slug'],
                       level=vehicle['entity']['roman_level'],
                       img=vehicle['entity']['image_url'],
                       current_count=vehicle['current_count'],
                       price=vehicle['price']['value'],
                       next_price=vehicle['next_price']['value'] if vehicle['next_price_timestamp'] else None
                       ) if vehicle['available'] else None

    async def poll_data(self) -> List[Vehicle]:
        logging.debug("start polling data")
        resp = requests.get(url=self.URL)

        logging.debug("start parsing data")
        data = resp.json()
        vehicles = [self.parse_vehicle_data(v) for v in data['results']]
        vehicles = [v for v in vehicles if v if v]

        self._last_polled_data = vehicles
        self._last_polled_at = int(time.time())
        return vehicles

    async def get_data(self) -> List[Vehicle]:
        if not self._last_polled_data or int(time.time()) - self._last_polled_at > self._cache_timeout:
            return await self.poll_data()

        return self._last_polled_data
