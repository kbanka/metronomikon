FROM python:3

WORKDIR /usr/src/app

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY tests ./tests
COPY run-test.py ./

ENTRYPOINT [ "python", "run-test.py" ]
