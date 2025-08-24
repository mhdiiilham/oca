from typing import List, Optional
from pydantic import BaseModel

class Item(BaseModel):
    name: str
    price: Optional[float] = None
    quantity: Optional[int] = None

class Receipt(BaseModel):
    merchant: Optional[str]
    date: Optional[str]
    total: Optional[float]
    items: List[Item] = []
    raw_text: str