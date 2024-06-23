FROM golang:alpine3.20

WORKDIR /root/server

RUN apk add openssl zsh g++ make py3-pip mandoc man-pages git vim neovim tmux

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -ldflags "-s -w" -o out/main

EXPOSE 3000/tcp

ENTRYPOINT ["./out/main"]
