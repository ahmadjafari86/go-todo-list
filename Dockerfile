# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go env -w GO111MODULE=on
RUN go mod download
COPY . .
RUN go build -o /bin/server ./cmd/server

# Run stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates curl
COPY --from=builder /bin/server /bin/server
WORKDIR /app
ENV PORT=8282
EXPOSE 8282
CMD ["/bin/server"]
