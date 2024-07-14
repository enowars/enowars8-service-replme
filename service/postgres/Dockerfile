FROM postgres:latest

RUN apt-get update; apt-get install -y --no-install-recommends pwgen

COPY docker-entrypoint-wrapped.sh /usr/local/bin

RUN chown root:root /usr/local/bin/docker-entrypoint-wrapped.sh
RUN chmod 500 /usr/local/bin/docker-entrypoint-wrapped.sh

ENTRYPOINT ["docker-entrypoint-wrapped.sh"]
