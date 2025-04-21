FROM golang:1.24.2-alpine AS build
RUN apk update && apk add --no-cache gcc musl-dev librdkafka-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -o main ./cmd/main.go
FROM alpine:latest

RUN apk update && apk add --no-cache librdkafka

WORKDIR /app

COPY --from=build /app/main .

EXPOSE 7770

CMD ["./main"]
