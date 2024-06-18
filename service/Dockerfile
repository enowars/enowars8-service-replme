FROM node:alpine
RUN apk add g++ make py3-pip go

WORKDIR /root

COPY ./frontend/package.json ./frontend/package-lock.json ./frontend/
RUN npm install --prefix frontend

COPY ./backend/go.mod ./backend/go.sum ./backend/
RUN go mod -C backend download

COPY ./frontend ./frontend
RUN make -C frontend

COPY ./backend ./backend
RUN go build -C backend -o out/main

COPY ./image-go ./image-go

EXPOSE 6969/tcp

ENTRYPOINT ["./backend/out/main", "-d", "./frontend/out", "-i", "./image-go", "-k", "/apikey/key.txt"]

