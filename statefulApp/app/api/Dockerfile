FROM golang:alpine3.20 AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
RUN GOARCH=amd64 go build -o /app/build/api /app/cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/build/api /app/api
COPY --from=builder /app/sql /app/sql

CMD ["/app/api"]

EXPOSE 8080