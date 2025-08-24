from abc import ABC, abstractmethod
from oca.domain.models import Receipt

class ReceiptParser(ABC):
    @abstractmethod
    def parse(self, text: str) -> Receipt:
        pass
