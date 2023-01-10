FROM openjdk:8-jre-alpine
MAINTAINER yiji <yiji@apche.org>

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

ADD target/xrpc-server-1.0-SNAPSHOT.jar /home/admin/release/xrpc-server-1.0-SNAPSHOT.jar

ADD start.sh /home/admin/release/start.sh

RUN chmod a+x /home/admin/release/start.sh

ENTRYPOINT ["/bin/sh", "/home/admin/release/start.sh"]
