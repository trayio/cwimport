FROM alpine:latest

COPY cwimport /bin/cwimport
COPY config.hcl /etc/config.hcl

RUN apk add --update ca-certificates

ENTRYPOINT ["/bin/cwimport"]
