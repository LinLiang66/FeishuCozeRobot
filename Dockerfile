FROM golang:1.212-alpine

WORKDIR /go/src/coze_robot

COPY . .
ENV GOPROXY=https://goproxy.cn
RUN echo http://mirrors.aliyun.com/alpine/v3.10/community/ > /etc/apk/repositories \
    && echo http://mirrors.aliyun.com/alpine/v3.10/main/ >> /etc/apk/repositories \
    && apk add tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo 'Asia/Shanghai' >/etc/timezone
RUN  go build -o coze_robot .

CMD ["./coze_robot"]

