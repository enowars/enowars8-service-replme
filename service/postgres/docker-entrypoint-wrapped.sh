#!/usr/bin/env bash

if [[ -e /root/data/secret ]]; then
    POSTGRES_PASSWORD="$(cat /root/data/secret)"
else
    POSTGRES_PASSWORD="$(pwgen --numerals --capitalize --remove-chars="'\\" -1 32)"
    echo -n "$POSTGRES_PASSWORD" > /root/data/secret
fi
export POSTGRES_PASSWORD

docker-entrypoint.sh postgres
