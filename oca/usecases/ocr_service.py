from oca.interfaces.ocr_engine import OcrEngine
from oca.interfaces.parser import ReceiptParser
from oca.domain.models import Receipt

class OcrService:
    def __init__(self, engine: OcrEngine, parser: ReceiptParser):
        self.engine = engine
        self.parser = parser

    def process(self, image_bytes: bytes) -> Receipt:
        text = self.engine.extract_text(image_bytes)
        receipt = self.parser.parse(text)
        return receipt

