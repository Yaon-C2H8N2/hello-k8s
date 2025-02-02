FROM golang:alpine3.20 AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
RUN GOARCH=amd64 go build -o /app/build/api ./cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/build/api /app/api

CMD ["/app/api"]

EXPOSE 8080