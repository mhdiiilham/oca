from abc import ABC, abstractmethod
from src.domain.models import Receipt

class ReceiptParser(ABC):
    @abstractmethod
    def parse(self, text: str) -> Receipt:
        pass
