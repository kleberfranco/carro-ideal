FROM golang:1.24-alpine AS dev
WORKDIR /app

ENV GO111MODULE=on
ENV PATH="/go/bin:${PATH}"

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/cosmtrek/air@v1.49.0

COPY . .

CMD ["air", "-c", ".air.toml"]


FROM golang:1.24-alpine AS builder
WORKDIR /app

ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./app/main.go


FROM alpine:latest AS prod
WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/web ./web

EXPOSE 8080

CMD ["./server"]