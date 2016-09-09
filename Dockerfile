FROM ubuntu:16.04

ENV TERRAFORM_VERSION=0.6.16
ENV DOCKER_COMPOSE_VERSION=1.7.1
ENV DOCKER_MACHINE_VERSION=v0.7.0
ENV RANCHER_CLI_VERSION=v0.0.1
ENV RANCHER_COMPOSE_VERSION=v0.8.6

RUN apt-get update && \
    apt-get install -y curl wget unzip apt-transport-https ca-certificates && \
    apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D && \
    echo "deb https://apt.dockerproject.org/repo ubuntu-xenial main" > /etc/apt/sources.list.d/docker.list && \
    apt-get update && \
    apt-get install -y docker-engine && \
    apt-get clean -q && rm -rf /var/lib/apt/lists/* && \
    wget -O /tmp/terraform.zip https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip -d /usr/local/bin /tmp/terraform.zip && \
    curl -L https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-Linux-x86_64 > /usr/local/bin/docker-compose && \
    chmod +x /usr/local/bin/docker-compose && \
    curl -L https://github.com/docker/machine/releases/download/${DOCKER_MACHINE_VERSION}/docker-machine-Linux-x86_64 > /usr/local/bin/docker-machine && \
    chmod +x /usr/local/bin/docker-machine && \
    curl -L https://github.com/rancher/cli/releases/download/${RANCHER_CLI_VERSION}/rancher-linux-amd64.tar.gz | \
        tar -xvz -C /tmp && \
    mv /tmp/rancher-${RANCHER_CLI_VERSION}/rancher /usr/local/bin && \
    curl -L https://releases.rancher.com/compose/${RANCHER_COMPOSE_VERSION}/rancher-compose-linux-amd64-${RANCHER_COMPOSE_VERSION}.tar.gz | \
        tar -xvz -C /tmp && \
    mv /tmp/rancher-compose-${RANCHER_COMPOSE_VERSION}/rancher-compose /usr/local/bin && \
    mkdir /odoko

WORKDIR /odoko
ADD /container /odoko
ADD bin/* /usr/local/bin/

ENTRYPOINT ["installer"]
