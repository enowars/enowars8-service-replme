FROM postgres@sha256:0aafd2ae7e6c391f39fb6b7621632d79f54068faebc726caf469e87bd1d301c0

RUN apt-get update; apt-get install -y --no-install-recommends pwgen

COPY docker-entrypoint-wrapped.sh /usr/local/bin

RUN chown root:root /usr/local/bin/docker-entrypoint-wrapped.sh
RUN chmod 500 /usr/local/bin/docker-entrypoint-wrapped.sh

ENTRYPOINT ["docker-entrypoint-wrapped.sh"]
