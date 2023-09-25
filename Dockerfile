FROM golang:1.21-alpine

ARG CERT_FILE=./CatoNetworksTrustedRootCA.cer
COPY ${CERT_FILE} /etc/ssl/certs/ca-certificates.crt
RUN apk --no-cache add ca-certificates
RUN cp /etc/ssl/certs/ca-certificates.crt /usr/local/share/ca-certificates
RUN update-ca-certificates
