#FROM golang:alpine
#COPY config/ /server/config/
#ADD ./http-server /server/http-server
#RUN  chmod u+x /server/http-server
#WORKDIR /server
#EXPOSE 50051
#COPY run-http.sh /run-http.sh
#RUN chmod u+x /run-http.sh
#CMD ["/run-http.sh"]

FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build
EXPOSE 50051
ADD go.mod .
ADD go.sum .

COPY ./submodule/services-proto.git ./submodule/services-proto.git
COPY ./submodule/support-go.git ./submodule/support-go.git

RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/http-server http-server.go

# nomal
#FROM alpine
# smaller
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app

COPY ./config /app/config

COPY --from=builder /app/http-server /app/http-server

CMD ["./http-server"]
