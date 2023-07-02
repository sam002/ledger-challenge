FROM python:3
COPY test/requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt