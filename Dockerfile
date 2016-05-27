FROM alpine:latest

COPY cwimport /bin/cwimport

RUN apk add --update ca-certificates

ENTRYPOINT ["/bin/cwimport"]
