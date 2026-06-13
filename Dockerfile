FROM golang:1.22-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o server ./app/main.go

FROM alpine:3.18 AS prod
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/web ./web
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./server"]
