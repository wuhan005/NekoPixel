FROM golang:1.20-alpine as go_builder

WORKDIR /app

ENV CGO_ENABLED=0

ARG GITHUB_SHA=dev

COPY . .

RUN go mod tidy
RUN go build -v -ldflags "-w -s -extldflags '-static' -X 'github.com/wuhan005/NekoPixel/internal/conf.BuildCommit=$GITHUB_SHA'" -o NekoPixel ./cmd/

FROM alpine:latest

RUN apk update && apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone

WORKDIR /home/app

COPY --from=go_builder /app/NekoPixel .

RUN chmod 777 /home/app/NekoPixel

ENTRYPOINT ["./NekoPixel"]
EXPOSE 8080
