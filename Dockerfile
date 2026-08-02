FROM alpine:3.23

COPY stamp /usr/local/bin/stamp

WORKDIR /app
CMD ["/usr/local/bin/stamp"]
