FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy && go build -o /app/user-srv ./

FROM alpine
WORKDIR /app
COPY --from=builder /app/user-srv /app/user-srv
RUN chmod +x /app/user-srv

ENTRYPOINT ["/app/user-srv"]