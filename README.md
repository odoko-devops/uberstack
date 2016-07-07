Odoko Development Environment Setup
===================================
Introduction
------------
The configurations in this project are intended to automate the creation of a 
containerised software development stack. This will give you:

 * Docker Registry (with SSL and basic auth)
 * Rancher Server (with basic auth)
 * Jenkins Server (with basic auth)

The basic steps should be:
 * Invoke Terraform
 * Wait for it to give you the IP address of your server
 * Configure this in your DNS
 * Wait for the Registry container to notice the DNS change, and configure itself

Usage
-----
Copy the terraform.tfvars.example file and name it terraform.tfvars. Replace the
placeholders there with sensible values.

To create the infrastructure, do:

    terraform apply

Once the node is up, Terraform will announce the IP address of the host. You will
need to point your domain name (e.g. docker.example.com) at this IP address. The
Registry container will eventually spot that this name has been set up, and will
proceed with creating SSL certificates, and then start up the Registry server.


Inspired by: https://www.airpair.com/aws/posts/ntiered-aws-docker-terraform-guide
