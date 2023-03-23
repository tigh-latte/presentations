FROM golang:1.20-buster as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bank ./cmd/bank

FROM alpine

COPY --from=builder /app/bank /bank

ENTRYPOINT [ "/bank" ]
