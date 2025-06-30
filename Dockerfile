# Go backend build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy && go build -o app ./cmd/root.go

# Final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY configs ./configs
EXPOSE 8080
CMD ["./app"]
