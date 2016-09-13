FROM ubuntu:16.04

RUN apt-get update && \
    apt-get install -y curl git && \
    mkdir /odoko && \
    cd /usr/local && \
    curl https://storage.googleapis.com/golang/go1.7.linux-amd64.tar.gz | tar -xz 

ENV GOPATH=/odoko/golibs
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

WORKDIR /odoko
RUN go get gopkg.in/yaml.v2 github.com/kr/pty

ENTRYPOINT ["./build.sh"]
