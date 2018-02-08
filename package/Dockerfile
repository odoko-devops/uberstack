FROM alpine:3.5
RUN apk add --no-cache ca-certificates openssh-client
COPY uber /usr/bin/
WORKDIR /mnt
ENTRYPOINT ["uber"]
CMD  ["--help"]
