import logging
from abc import ABC, abstractmethod
from dataclasses import dataclass
from typing import List


@dataclass
class MonitoringData:
    vehicle_id: str
    chat_id: int
    minimal_count: int


class AbstractStorage(ABC):
    @abstractmethod
    async def save(self, data: MonitoringData):
        pass

    @abstractmethod
    async def remove(self, chat_id: int, vehicle_id: str):
        pass

    @abstractmethod
    async def get_by_vehicle_id(self, vehicle_id) -> List[MonitoringData]:
        pass

    @abstractmethod
    async def get_all(self) -> List[MonitoringData]:
        pass


class InMemoryStorage(AbstractStorage):
    def __init__(self):
        self._data: List[MonitoringData] = []

    async def save(self, data: MonitoringData):
        try:
            existent = next(s for s in self._data if s.chat_id == data.chat_id and s.vehicle_id == data.vehicle_id)
            self._data.remove(existent)
        except StopIteration:
            # object does not exist in data
            pass
        finally:
            self._data.append(data)

    async def remove(self, chat_id: int, vehicle_id: str):
        try:
            existent = next(s for s in self._data if s.chat_id == chat_id and s.vehicle_id == vehicle_id)
            self._data.remove(existent)
        except StopIteration:
            logging.debug(f'data with chat_id:{chat_id} and vehicle_id: {vehicle_id} does not exist')

    async def get_by_vehicle_id(self, vehicle_id) -> List[MonitoringData]:
        return [s for s in self._data if s.vehicle_id == vehicle_id]

    async def get_all(self) -> List[MonitoringData]:
        return self._data
