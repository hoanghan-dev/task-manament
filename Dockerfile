FROM golang:1.26.3-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG APP=api

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/${APP}

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./app"]