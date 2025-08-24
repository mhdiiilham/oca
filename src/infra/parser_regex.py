import re
from src.domain.models import Receipt, Item
from src.interfaces.parser import ReceiptParser

class RegexReceiptParser(ReceiptParser):
    def parse(self, text: str) -> Receipt:
        total = None
        match = re.search(r"total[:\s]*([\d,.]+)", text, re.IGNORECASE)
        if match:
            total = float(match.group(1).replace(",", "").replace(".", "")) / 100

        return Receipt(
            merchant=None,
            date=None,
            total=total,
            items=[],
            raw_text=text
        )
