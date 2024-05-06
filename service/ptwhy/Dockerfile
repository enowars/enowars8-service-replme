FROM node:alpine3.19
RUN apk add g++ make py3-pip openssl zsh mandoc man-pages
COPY . .
RUN rm -rf node_modules
RUN rm -rf out
RUN npm i
RUN npm run build
EXPOSE 3000/tcp
ENTRYPOINT ["node", "out/server.js"]

