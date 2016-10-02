FROM golang:1.7.1

RUN go get  gopkg.in/yaml.v2 \
            github.com/kr/pty

ENV USER root
WORKDIR /go

ENTRYPOINT ["./build.sh"]
