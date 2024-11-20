# Build stage
FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bankapp main.go

# Run stage
FROM scratch

COPY --from=builder /app/bankapp /bankapp

EXPOSE 8080

ENTRYPOINT ["/bankapp"]
