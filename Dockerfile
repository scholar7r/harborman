# -- Building stage ----------
FROM golang:1.27 AS builder
WORKDIR /workspace

RUN go install github.com/go-task/task/v3/cmd/task@v3

COPY . .
RUN task build

# -- Production stage ----------
FROM debian:bookworm-slim AS production
WORKDIR /workspace

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /workspace/build .
