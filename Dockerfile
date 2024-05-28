FROM golang:1.22.3-alpine as builder
WORKDIR /backend-app-files

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./application ./cmd/main.go

FROM alpine

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /backend-app-files/application /application

ENTRYPOINT ["/application"]