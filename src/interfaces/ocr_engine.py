from abc import ABC, abstractmethod
from src.domain.models import Receipt

class OcrEngine(ABC):
    @abstractmethod
    def extract_text(self, image_bytes: bytes) -> str:
        pass
