FROM node:alpine AS frontend

WORKDIR /app

COPY ./frontend/package.json .

RUN yarn install

COPY ./frontend/ .

RUN ls -laR .

RUN yarn build

FROM golang:1.17-alpine AS backend

WORKDIR /app

COPY . .

RUN go build

FROM alpine AS app

WORKDIR /app

COPY --from=backend /app/migadu-admin .

COPY --from=frontend /app/build ./frontend

EXPOSE 5000

CMD ./migadu-admin
