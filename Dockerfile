FROM golang:1.24-alpine AS builder

ENV TZ Asia/Shanghai

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
COPY . .

RUN go build -ldflags="-s -w" -o /app/user-srv .

FROM alpine:3.20

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/user-srv /app/user-srv
RUN chmod +x /app/user-srv

ENTRYPOINT ["/app/user-srv"]