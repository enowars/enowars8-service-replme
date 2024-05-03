FROM node:alpine3.19
RUN apk add g++ make py3-pip
COPY package* .
RUN npm i
COPY server.js .
EXPOSE 3000/tcp
ENTRYPOINT ["node", "server.js"]
# ENTRYPOINT ["/bin/sh" "server.js"]

