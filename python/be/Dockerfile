FROM python:3.9-slim
ARG VERSION
COPY python/be/requirements.txt /
RUN pip3 install -r /requirements.txt
RUN python -m nltk.downloader punkt
COPY python/be/app /app
WORKDIR /app
CMD ["gunicorn", "-c", "/app/gunicorn.conf.py", "app:app"]