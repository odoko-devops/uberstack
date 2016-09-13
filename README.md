Uberstack
=========

Introduction
------------
Uberstack is a tool that simplifies working with Docker and Rancher. It has
two main focuses - to simplify the creation of infrastructure, and to 
simplify the deployment of applications.

Infrastructure comes in a number of guises. If you are deploying, for example,
at Amazon Web Services, then constructing a VPC, and establishing security
groups would be necessary preparations that can be automated.

Then, to effectively use Rancher, as well as a functioning RAncher server,
we also need some other supporting infrastructure, such as a private Docker 
Registry, a Jenkins server and perhaps a VPN server.

Once we have Rancher and related infrastructure in place, we can deploy
hosts onto which we can push our applications.

Uberstack can help with all of this. Thus, it can:

 * Install a management host in a local Virtualbox, with Rancher Server
 * Install a local node for Docker development, connected to Rancher
 * Create your VPC at Amazon
 * Create a management server, which hosts a:
   * Rancher server
   * Private Docker registry
   * Jenkins server
   * VPN server
 * Create hosts for Docker usage, connect them to Rancher with predefined
   labels
 * Deploy multiple application stacks to Docker/Rancher, supporting 
   multiple 'stacks' configured for multiple 'environments'

The last feature is sufficient to use Uberstack, without requiring any of the
previous steps.

All of the above is configured with simple YAML configuration files.

Installing and Initialising Uberstack
-------------------------------------
Download the latest release of Uberstack from 
https://github.com/odoko-devops/uberstack/releases.

You will require both files (uberstack and uberstack-remote-agent), and should 
place both of them on your path.

Uberstack requires one or both of two environment variables:

 * UBER_HOME: the location where 'stacks' and 'uberstacks' are defined, 
   including docker-compose.yml and rancher-compose.yml files for building 
   application stacks.
 * UBER_STATE: where Uberstack looks for certain configuration and certain 
   state information, particularly when building hosts and related 
   infrastructure.

With these environment variables defined, initialise Uberstack with:

    uberstack init
    
This will download dependency binaries into your UBER_STATE directory for 
future use.

Configuring Uberstack
---------------------
Uberstack is configured with a `config.yml` file inside your UBER_STATE 
directory. Uberstack will, soon, have a mechanism to provide a sample config
file to start with.

This configuration file defines infrastructure and hosts. Applications have
their own YAML configuration files.

Deploying Infrastructure
------------------------
Uberstack supports multiple 'providers' in much the same way as docker-machine
(which is used under the bonet). At present, it supports `amazonec2` and 
`virtualbox`, which gives support for local, and remote development and 
deployment. Other providers can be supported.

The following will initialise the AmazonEC2 VPN as configured:

    uberstack provider up amazonec2

 Note: This functionality is currently broken due to refactoring. It will be
 fixed shortly, and documentation about how to configure it will be provided.

Deploying Hosts
---------------
Hosts can be created, destroyed or replaced with Uberstack:

    uberstack host up <hostname>
    uberstack host rm <hostname>
    uberstack host replace <hostname>

The `<hostname>` mentioned above is a reference to a host configured in
`config.yml`.

Deploying Applications
----------------------
Uberstack for Applications has a different YAML based configuration setup.

It shares the state.yml file with the other parts of Uberstack for recording
Rancher credentials and provider details, but otherwise, its configuration is
separate.

The name "Uberstack" comes from the idea of building applications out of a
stack of stacks. Docker and Rancher suggest the idea of a 'stack' which can
be expressed with a `docker-compose.yml` file, and enhanced with a 
`rancher-compose.yml` one. An "Uberstack" groups these stacks together,
and providing configuration (environment variables) for one or more 
'environments', such as local, dev, staging, qa and production.


