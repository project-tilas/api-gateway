FROM alpine:3.7

EXPOSE 8080

ADD api-gateway /bin/api-gateway

ENTRYPOINT "api-gateway"