FROM golang:1.22.3-alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./app ./cmd/main.go

FROM alpine:latest

COPY --from=builder /app/app /app

ENTRYPOINT ["/app"]