FROM ubuntu:16.04

ENV TERRAFORM_VERSION=0.6.16
ENV DOCKER_COMPOSE_VERSION=1.7.1
ENV DOCKER_MACHINE_VERSION=v0.7.0
ENV RANCHER_CLI_VERSION=v0.0.1

RUN apt-get update && \
    apt-get install -y curl wget unzip python-yaml apt-transport-https ca-certificates python-pexpect && \
    apt-get install -y python-pip python-dev dh-autoreconf && \
    apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D && \
    echo "deb https://apt.dockerproject.org/repo ubuntu-xenial main" > /etc/apt/sources.list.d/docker.list && \
    apt-get update && \
    apt-get install -y docker-engine && \
    apt-get clean -q && rm -rf /var/lib/apt/lists/* && \
    pip install --upgrade pip && \
    pip install requests pyjq && \
    wget -O /tmp/terraform.zip https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip -d /usr/local/bin /tmp/terraform.zip && \
    curl -L https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-Linux-x86_64 > /usr/local/bin/docker-compose && \
    chmod +x /usr/local/bin/docker-compose && \
    curl -L https://github.com/docker/machine/releases/download/${DOCKER_MACHINE_VERSION}/docker-machine-Linux-x86_64 > /usr/local/bin/docker-machine && \
    chmod +x /usr/local/bin/docker-machine && \
    curl -L https://github.com/rancher/cli/releases/download/${RANCHER_CLI_VERSION}/rancher-linux-amd64.tar.gz | \
    tar -xvz -C /tmp && \
    mv /tmp/rancher-${RANCHER_CLI_VERSION}/rancher /usr/local/bin && \
    mkdir /odoko

ADD . /odoko
WORKDIR /odoko

ENTRYPOINT [ "python", "-u", "main.py"]
CMD ["up"]
