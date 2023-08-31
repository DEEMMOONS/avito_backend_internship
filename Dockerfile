FROM golang:1.18-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o main cmd/main.go

FROM alpine:latest

COPY --from=builder . .
EXPOSE 80

CMD ["app/main"]
