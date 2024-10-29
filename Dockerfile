# Step 1: Build from golang alpine
FROM golang:1.23-alpine as builder
WORKDIR /go/src/app
COPY cmd ./cmd
COPY internal ./internal
COPY go.mod go.sum ./
COPY vendor ./vendor

# Build the binary
ARG GO_BINARY
ENV GO_BINARY ${GO_BINARY:-./cmd/bot}
RUN CGO_ENABLED=0 go build -mod vendor -buildvcs=false -o /go/bin/app ${GO_BINARY}

# Step 2: Copy the binary to a new container
FROM alpine:latest
COPY --from=builder /go/bin/app /app
ENTRYPOINT /app

