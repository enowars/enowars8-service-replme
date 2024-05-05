FROM node:alpine3.19
RUN apk add g++ make py3-pip
COPY . .
RUN npm i
EXPOSE 3000/tcp
ENTRYPOINT ["node", "server.js"]

