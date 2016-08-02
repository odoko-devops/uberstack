FROM ubuntu:16.04

RUN apt-get update && \
    apt-get install -y python-pip python-dev dh-autoreconf && \
    apt-get clean -q && \
    pip install --upgrade pip && \
    pip install requests pyjq && \
    rm -rf  /tmp/* /var/tmp/* /var/lib/apt/lists/*

ADD /start.py /

CMD ["/start.py"]
