FROM alpine:3.22

COPY stamp /usr/local/bin/stamp

WORKDIR /app
CMD ["/usr/local/bin/stamp"]
