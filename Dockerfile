FROM alpine:3.22

COPY stamp /usr/local/bin/stamp

CMD ["/usr/local/bin/stamp"]
