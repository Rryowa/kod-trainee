FROM golang:alpine AS builder
RUN apk --no-cache add bash git make
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/main ./cmd/main.go

ENTRYPOINT ["./bin/main"]