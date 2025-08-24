import pytesseract
from PIL import Image
import io
from src.interfaces.ocr_engine import OcrEngine

class TesseractOcrEngine(OcrEngine):
    def extract_text(self, image_bytes: bytes) -> str:
        img = Image.open(io.BytesIO(image_bytes))
        return pytesseract.image_to_string(img)
