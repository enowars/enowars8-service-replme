FROM node:alpine3.19
RUN apk add g++ make py3-pip openssl zsh mandoc man-pages

WORKDIR /root/server

COPY . .

RUN npm install
RUN npm run build

EXPOSE 3000/tcp

ENTRYPOINT ["node", "out/main.js"]

