FROM nginx

RUN apt-get update && \
    apt-get install -y apache2-utils

ADD /nginx.conf /etc/nginx-configuration

ADD start.sh /
CMD ["bash", "/start.sh"]

