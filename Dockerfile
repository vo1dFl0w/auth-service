FROM golang:1.24.3-alpine AS builder

WORKDIR /auth-service

RUN apk --no-cache add git bash make gcc gettext musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

ENV CONFIG_PATH=configs/config.yaml
ENV CGO_ENABLED=0

RUN go build --ldflags="-w -s" -o auth-service ./cmd/auth-service

FROM alpine AS runner
RUN apk add --no-cache ca-certificates

WORKDIR /auth-service

COPY --from=builder /auth-service/configs /auth-service/configs
COPY --from=builder /auth-service/auth-service /auth-service/auth-service

CMD ["./auth-service"]

