FROM python:3.11.9@sha256:c5254471e6073d8942091227f469302f85b30d4b23077226f135360491f5226a

RUN useradd -ms /bin/bash -u 1000 checker

RUN apt update -y

USER checker
WORKDIR /checker

COPY ./pyproject.toml ./poetry.lock .
RUN pip install poetry
RUN python -m poetry install

COPY ./src .

ENTRYPOINT [ "python", "-m", "poetry", "run", "gunicorn", "-c", "gunicorn.conf.py", "checker:app" ]

