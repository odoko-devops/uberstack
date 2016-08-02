FROM jenkins

USER root
RUN apt-get update && \
    apt-get install -y git curl zip maven && \
    rm -rf /var/lib/apt/lists/*

RUN apt-get update && \
    apt-get install -y apt-transport-https ca-certificates && \
    apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D && \
    echo "deb https://apt.dockerproject.org/repo ubuntu-trusty main" > /etc/apt/sources.list.d/docker.list && \
    apt-get update && \
    apt-get install -y docker-engine && \
    gpasswd -a jenkins docker

RUN mkdir /plugins && \
    curl -s http://mirrors.jenkins-ci.org/plugins/git/latest/git.hpi > /plugins/git.hpi

USER ${user}

#ADD /start.sh /start.sh

#ENTRYPOINT [ "bash", "/start.sh" ]
