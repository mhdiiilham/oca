# OCA (OCR API Service)

OCA is a Python-based OCR (Optical Character Recognition) microservice designed to extract text from receipt images.
It is intended to work as a backend service that integrates with other systems (e.g., a Go service or Telegram bot).

## Features
- REST API with FastAPI
- OCR using Tesseract
- Input validation with Pydantic
- Clean architecture folder structure
- Easy to integrate with other services

## Requirements
- Python 3.9+
- Tesseract OCR installed on your system
- Virtual environment (recommended)

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/oca.git
   cd oca
   ```

2. Create and activate virtual environment:
   ```bash
   python3 -m venv venv
   source venv/bin/activate   # On macOS/Linux
   venv\Scripts\activate    # On Windows
   ```

3. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```

4. Run the service:
   ```bash
   uvicorn app.main:app --reload
   ```

The service will be available at: [http://localhost:8000](http://localhost:8000)

## Example API Usage

### Upload a Receipt
```bash
curl -X POST "http://localhost:8000/ocr" -F "file=@receipt.jpg"
```

Response:
```json
{
  "text": "Extracted receipt content here..."
}
```

## Project Structure
```
oca/
│── app/
│   ├── main.py          # FastAPI entrypoint
│   ├── api/             # API routes
│   ├── core/            # Core logic
│   ├── models/          # Pydantic models
│   ├── services/        # OCR service
│   └── utils/           # Helper functions
│── tests/               # Unit tests
│── requirements.txt     # Python dependencies
│── README.md            # Documentation
│── .gitignore
```

## License
This project is licensed under the [MIT License](./LICENSE).
