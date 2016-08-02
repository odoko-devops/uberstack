Odoko Development Environment Setup
===================================
Introduction
------------
The configurations in this project are intended to automate the creation of a 
containerised software development stack, at AWS. This will give you:

 * Docker Registry (with SSL and basic auth)
 * Rancher Server (with basic auth)
 * Jenkins Server (with basic auth)

The basic steps should be:
 * Request an elastic IP via the AWS console
 * Assign three domain names to this IP, one for each of docker, rancher and
   jenkins
 * Edit config.yml accordingly
 * Run the Odoko Stack docker container
 * Wait

After some minutes, you should have a fully functioning server.

Usage
-----
Copy the config.yml.example file to config.yml, and replace the
placeholders there with sensible values.

To create the infrastructure, do:

    bin/make-image
    bin/build-me-a-server

To create a docker host, edit the config.yml setting .hosts.count to the
number of nodes that you want to have, then:

    bin/create-my-docker-nodes
