FROM registry:2.4.1

RUN apt-get update && \
    apt-get install -y cron git curl apt-utils && \
    git clone https://github.com/certbot/certbot.git /certbot && \
    /certbot/letsencrypt-auto --os-packages-only && \
    echo "30 2 * * 1 root /certbot/letsencrypt-auto renew >> /var/log/letsencrypt-renew.log" > /etc/cron.d/letsencrypt && \
    sed -i "s/addr: :5000/addr: :443/" /etc/docker/registry/config.yml

ADD /startup.sh /

ENTRYPOINT [ "/startup.sh" ]

CMD [ "run" ]
