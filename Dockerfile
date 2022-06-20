# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY . ./
RUN go mod download

RUN go build -o /pack-and-go ./app

EXPOSE 8080

CMD [ "/pack-and-go" ]