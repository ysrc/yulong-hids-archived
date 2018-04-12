FROM golang:1.10 as builder

RUN mkdir -p /go/src/yulong-hids
WORKDIR /go/src/yulong-hids
ADD . /go/src/yulong-hids
RUN go build -o ./web/web --ldflags='-w -s -linkmode external -extldflags "-static"' ./web/main.go

FROM alpine
MAINTAINER Jason Cooper "mrderek@protonmail.com"
RUN mkdir /web
WORKDIR /web
COPY --from=builder /go/src/yulong-hids/web/ .
RUN cp ./conf/app-config-sample.conf ./conf/app.conf
RUN apk update
RUN apk upgrade
RUN apk add ca-certificates && update-ca-certificates
RUN apk add openssl
RUN apk add --update tzdata
ENV TZ=Asia/Shanghai
RUN rm -rf /var/cache/apk/*
RUN chmod +x /web/web
ENTRYPOINT [ "./web" ]
