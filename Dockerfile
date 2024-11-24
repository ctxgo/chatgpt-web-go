FROM  golang:1.21.10-alpine3.20 as builder

WORKDIR /go/release

COPY . .

RUN set -x \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk add gcc libc-dev \
    && go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod tidy \
    && CGO_ENABLED=1 \ 
    GOOS=linux GOARCH=amd64 \
    go build -o /ai-go \
    -x -ldflags="-w -s" \
    cmd/server/main.go 


FROM alpine:3.20
RUN  sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
   && apk add  --no-cache tzdata 
COPY --from=builder /ai-go /
CMD ["/ai-go","-c","/config.yaml"]
