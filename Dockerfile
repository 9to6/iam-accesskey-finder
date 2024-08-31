FROM golang:1.22-alpine as builder
LABEL authors="gunner"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

ADD . /app
RUN go build -o ak

FROM alpine:3

COPY --from=builder /app/ak /bin/ak

ENTRYPOINT ["ak"]

