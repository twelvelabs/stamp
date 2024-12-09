FROM alpine:3.21

COPY stamp /usr/local/bin/stamp

CMD ["/usr/local/bin/stamp"]
