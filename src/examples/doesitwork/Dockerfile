FROM node:alpine

WORKDIR /app

COPY package.json yarn.lock /app/

RUN yarn install

COPY index.js /app/

CMD yarn start
