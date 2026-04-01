FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/api ./cmd/api

FROM alpine:3.21

RUN addgroup -S app && adduser -S app -G app
WORKDIR /app

COPY --from=builder /bin/api /app/api

EXPOSE 8080
USER app

ENTRYPOINT ["/app/api"]
