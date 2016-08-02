Full Docker Registry
====================
This project aims to make the installation of a Docker Registry, with
basic authentication and SSL (with letsencrypt.org), trivial.

Usage
-----
This repo creates a Docker image that will run a Docker Registry.

It expects the following environment variables:

 * DOMAIN: the domain on which the Registry will serve
 * EMAIL: the email address associated with the Registry
 * USERNAME: a single username with which to secure the Registry
 * PASSWORD: a single password to secure the above username

When the above are provided and an container is started, first the
username/password are added to the htpasswd file for authentication,
then a letsencrypt SSL certificate is requested. For this to work,
the domain ($DOMAIN) must be pointing to the host on which the
container is running, and port 443 must be accessible from the 
outside (as letsencrypt will attempt to access the server whilst
generating the SSL certificate).

The container will wait for DNS to be updated before generating the
cert and starting the Registry.

The container can be run with CMD=passwd and another USERNAME/PASSWORD
combination to add additional users to the Registry.

When run with CMD=cron, it will start a cron task that will periodically
renew the SSL certificate.
