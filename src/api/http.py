from fastapi import FastAPI, UploadFile, File
from src.usecases.ocr_service import OcrService
from src.infra.tesseract_engine import TesseractOcrEngine
from src.infra.parser_regex import RegexReceiptParser

app = FastAPI()

ocr_service = OcrService(
    engine=TesseractOcrEngine(),
    parser=RegexReceiptParser()
)

@app.post("/scan")
async def scan_receipt(file: UploadFile = File(...)):
    image_bytes = await file.read()
    receipt = ocr_service.process(image_bytes)
    return receipt.dict()
