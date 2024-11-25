FROM python:3.9-slim
WORKDIR /app
COPY . /app

RUN pip install --no-cache-dir -r req.txt

EXPOSE 50051

ENV NAME World

CMD ["python", "main.py"]