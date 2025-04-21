FROM golang:1.24.2 AS build

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc libc6-dev librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags="-w -s" -o main ./cmd/main.go

FROM debian:bookworm-slim AS final

RUN apt-get update && apt-get install -y --no-install-recommends \
    librdkafka1 ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /app/main .

EXPOSE 7770

CMD ["./main"]
