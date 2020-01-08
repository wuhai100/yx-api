FROM alpine:latest

RUN apk add --update ca-certificates
RUN update-ca-certificates
RUN apk add --update tzdata
ENV TZ=Asia/Shanghai
COPY yx-api /
RUN  mkdir /public
RUN  mkdir /swagger
COPY swagger /swagger
RUN mkdir -p /usr/local/go/lib/time/
COPY zoneinfo.zip /usr/local/go/lib/time/
RUN chmod +x /yx-api

CMD ["/yx-api"]
