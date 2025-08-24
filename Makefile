.PHONY: venv install run dev format lint clean

# Create virtual environment
venv:
	python3 -m venv venv

# Install dependencies
install:
	venv/bin/pip install -r requirements.txt

# Run the API service
run:
	venv/bin/python cmd/main.py

# Run with uvicorn (development mode, hot reload)
dev:
	venv/bin/uvicorn cmd.main:app --reload

# Format code with black
format:
	venv/bin/black .

# Lint code with flake8
lint:
	venv/bin/flake8 .

# Clean up venv and cache files
clean:
	rm -rf venv __pycache__ .pytest_cache .mypy_cache
