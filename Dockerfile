FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.24.1
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/app ./cmd/app/main.go
FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /bin/app /app/app
COPY --from=builder /go/bin/goose /app/goose
COPY --from=builder /app/migrations /app/migrations
WORKDIR /app
EXPOSE 8081
CMD ["/app/app"]