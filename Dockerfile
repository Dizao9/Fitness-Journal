FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/app ./cmd/app/main.go
FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /bin/app /app/app
WORKDIR /app
EXPOSE 8081
CMD ["/app/app"]