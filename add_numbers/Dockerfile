FROM python:3.11-alpine3.18

WORKDIR /app

COPY math_operations.py .
COPY test_math_operations.py .

CMD ["python", "-m", "unittest", "test_math_operations.py"]