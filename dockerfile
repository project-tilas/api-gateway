FROM alpine:3.7

EXPOSE 8081

ADD api-gateway /bin/api-gateway

ENTRYPOINT "api-gateway"