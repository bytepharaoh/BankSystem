FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o simplebank main.go

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/simplebank /app/simplebank
COPY app.env /app/app.env

EXPOSE 8080
CMD ["/app/simplebank"]
