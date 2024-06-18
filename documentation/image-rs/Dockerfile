FROM rust:alpine3.20
RUN apk add g++ make py3-pip openssl zsh mandoc man-pages git

WORKDIR /root/server

COPY . .

RUN cargo install --path .

EXPOSE 3000/tcp

ENTRYPOINT ["image-rs"]

