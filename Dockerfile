FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o simplebank main.go

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates netcat-openbsd

COPY --from=builder /app/simplebank /app/simplebank
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY app.env /app/app.env
COPY db/migration /app/db/migration
COPY scripts/docker-entrypoint.sh /app/docker-entrypoint.sh

RUN chmod +x /app/docker-entrypoint.sh

EXPOSE 8080
ENTRYPOINT ["/app/docker-entrypoint.sh"]
