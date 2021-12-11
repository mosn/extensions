FROM golang:1.14.13

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

RUN go get github.com/go-delve/delve/cmd/dlv@v1.6.1

# install tools.
RUN apt update && apt -y install net-tools sudo openbsd-inetd telnet

# change container root pwd
RUN  useradd -ms /bin/bash admin && \
     echo "root:root" | chpasswd && \
     echo "admin:admin" | chpasswd

WORKDIR /go