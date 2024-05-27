FROM alpine:3.20

COPY stamp /usr/local/bin/stamp

CMD ["/usr/local/bin/stamp"]
