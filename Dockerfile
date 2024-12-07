FROM golang:1.23-alpine as builder

RUN apk add --no-cache git build-base

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o auth cmd/main.go

FROM alpine:3.18

RUN apk add --no-cache ca-certificates

RUN wget https://golang.org/dl/go1.23.0.linux-amd64.tar.gz && \
    tar -C /usr/local -xvzf go1.23.0.linux-amd64.tar.gz && \
    rm go1.23.0.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app

COPY --from=builder /app/auth /usr/local/bin/auth

COPY --from=builder /app/.env /app/.env 
COPY --from=builder /app/storage/ /app/storage/
COPY --from=builder /app/internal /app/internal
COPY --from=builder /app/go.mod /app/go.mod
COPY --from=builder /app/go.sum /app/go.sum

RUN chmod 644 /app/.env

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/auth"]

