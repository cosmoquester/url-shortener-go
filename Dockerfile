FROM golang:1.16.5-alpine AS builder

# Set necessary environmet variables needed for running on scratch
ENV CGO_ENABLED=0

WORKDIR /build

# Download Packages
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

# Build
RUN go build -o url-shortener-go .

FROM scratch

WORKDIR /app

# Use only compiled binary
COPY --from=builder /build/url-shortener-go .

ENTRYPOINT [ "./url-shortener-go" ]
