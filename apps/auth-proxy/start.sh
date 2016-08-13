#!/bin/sh

# Expand variables inside nginx-configuration file:
eval "echo \"$(cat /etc/nginx-configuration)\"" > /etc/nginx/conf.d/default.conf

nginx -g "daemon off;"
