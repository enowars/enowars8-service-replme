FROM node:alpine
RUN apk add g++ make py3-pip go

WORKDIR /app

COPY . .

# RUN rm -rf node_modules/ tmp/ ptwhy/node_modules/ ptwhy/out/

RUN npm i
RUN npm run build
RUN npm run tailwind-build
RUN go run github.com/a-h/templ/cmd/templ@latest generate -v
RUN go build .

EXPOSE 6969/tcp

ENTRYPOINT ["./cafedodo"]

