Odoko Docker Stack
==================

Introduction
------------
Docker is a powerfully transformative tool for modern software development.
However, to start using it, we need a number of services available to us,
and these can take undue time and effort to set up and configure correctly.

Rancher
-------
Rancher (http://www.rancher.com) is an orchestration platform that greatly
eases the management of suites of applications all running as Docker 
containers across a range of hosts. It provides both a powerful command
line for deploying and managing applications and a simple yet powerful
UI that allows you to manage and understand the behaviour of your 
infrastructure.

This Project
------------
This project includes code to automate the installation of a range of
components; currently:

 * A Docker Registry (SSL enabled and password protected)
 * A Rancher Server (with basic auth - for orchestrating containers/hosts)
 * A Jenkins Server (with basic auth - for managing builds)
 
Over time, additional applications will be added to this infrastructure.

Currently, this code only deploys to Amazon Web Services, into a VPC that it 
will create for you. It could easily be extended to support other hosting
providers.

Configuration
-------------
Configuration is handled within config.yml. Copy the sample config.yml.example
to config.yml and then edit the file as desired.

Usage
-----
This application requires Docker to be installed locally.

To start, you will need to create a Docker image:

    bin/make-image

Before starting your management server, ensure that an elastic IP address has
been created manually via the AWS console, the domain names for Docker,
Rancher and Jenkins have been pointed at this IP, and that the eipalloc value
for this IP, and the IP itself, have been added to the config.yml file. The
management server will not be fully functional until the domain names can be
resolved.

To start up a management server:

    bin/build-me-a-server

To start one or more Rancher/Docker nodes, do:

    bin/create-my-docker-nodes
