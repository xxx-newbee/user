FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o user-srv ./user

FROM alpine
WORKDIR /app
COPY --from=builder /app/user-srv .
COPY --from=builder /app/user/etc /app/etc
CMD ["./user-srv", "-f", "etc/user.yaml"]