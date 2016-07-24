FROM ubuntu:14.04

RUN apt-get update && \
    apt-get install -y curl jq && \
    apt-get clean -q && \
    rm -rf  /tmp/* /var/tmp/* /var/lib/apt/lists/*

ADD /start.sh /

CMD ["bash", "/start.sh"]
