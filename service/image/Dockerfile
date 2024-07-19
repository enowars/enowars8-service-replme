FROM golang:alpine3.20@sha256:8c9183f715b0b4eca05b8b3dbf59766aaedb41ec07477b132ee2891ac0110a07

WORKDIR /root/server

RUN apk add openssl zsh g++ make py3-pip mandoc man-pages git vim neovim tmux

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -ldflags "-s -w" -o out/main

EXPOSE 3000/tcp

ENTRYPOINT ["./out/main"]
