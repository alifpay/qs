FROM golang:alpine AS builder

RUN apk update && apk add --no-cache ca-certificates

# Move to working directory /build
WORKDIR /build

COPY go.mod .

RUN go mod tidy

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy our static executable
COPY --from=builder /build/main /main

# Run the api binary.
ENTRYPOINT ["/main"]