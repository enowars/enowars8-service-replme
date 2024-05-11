FROM node:alpine
RUN apk add g++ make py3-pip go

WORKDIR /root

COPY . .

RUN npm install --prefix frontend
RUN make -C frontend
RUN go build -C backend -o out/main

EXPOSE 6969/tcp

ENTRYPOINT ["./backend/out/main", "-d", "./frontend/out", "-i", "./image"]

