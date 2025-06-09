FROM golang:1.24-alpine AS base

WORKDIR /app

COPY pkg/ ./pkg/
COPY go.mod go.sum ./

# WORKDIR /app/collector-service
RUN go mod download

COPY . .

RUN go build -o collector-service app/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=base /app/collector-service .
# COPY --from=base /app/config.toml .

# ENV APP_CONFIG=config

EXPOSE 8080

CMD [ "./collector-service" ]