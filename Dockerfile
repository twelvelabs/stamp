FROM alpine:3.19

COPY stamp /usr/local/bin/stamp

CMD ["/usr/local/bin/stamp"]
