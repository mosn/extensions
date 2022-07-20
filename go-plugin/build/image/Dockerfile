FROM golang:1.17.9
MAINTAINER yiji@apache.org

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

RUN go install github.com/go-delve/delve/cmd/dlv@v1.7.3

# install tools.
RUN apt update && apt -y install net-tools sudo openbsd-inetd telnet tcpdump zip iptables neovim sysstat iproute2 iftop

# change container root pwd
RUN  useradd -ms /bin/bash admin && \
     echo "root:root" | chpasswd && \
     echo "admin:admin" | chpasswd

WORKDIR /go

# how to use this Dockefile ?
# docker buildx create --use --name=builder --driver docker-container
# docker buildx build --platform linux/arm64,linux/amd64,linux/386 -t zonghaishang/delve:v1.7.3 . --push
# docker buildx rm builder

# how to use this Dockefile build local image for development ?
# docker build -t zonghaishang/delve:v1.7.3 .
