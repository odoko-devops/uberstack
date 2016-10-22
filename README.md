Uberstack
=========

Introduction
------------
Uberstack is a tool that simplifies working with Docker and Rancher. It 
facilitates the creation of infrastructure, and simplifies the deployment of 
applications.

Uberstack brings the concept of 'providers', which take two flavours:

 * A host provider is used to create infrastructure. It may, for example,
   create an AWS VPC as well as bringing up hosts onto which services can
   be installed. The default host provider wraps Terraform, so can 
   work with any hosting provider that Terraform supports
 * An app provider is responsible for starting/stopping applications

Uberstack aims to make the task of using Docker easier right from the get-go.
Therefore, one of the example configurations provided will install a suite
of development tools, assuming nothing much more than an existing Amazon
Web Services account. This example setup will install:

 * A Docker Registry (with SSL and basic auth)
 * A Rancher Server (with basic auth and logged into Docker Registry)
 * A Jenkins Server (with basic auth and logged into Docker Registry)
 * A proxy server (so that the above will all work with nice domain names)

See `/examples/dev-stack/README.md` for more information on how to use this.

Uberstsack is entirely configured by (relatively) simple YAML files.

Installing and Initialising Uberstack
-------------------------------------
Download the latest release of Uberstack from 
https://github.com/odoko-devops/uberstack/releases.

Place the `uberstack` file onto your path.

Uberstack requires the `UBER_HOME` environment variable to be set. This 
defines where your YAML configuration files can be found.

Uberstack depends on a number of other tools, e.g. Rancher support requires
rancher-compose, and of course, Docker. At present, these must be present on
the host system, although Uberstack will soon have the ability to download
its own copies of these tools.

Configuring Uberstack
---------------------
Uberstack is configured with a set of YAML files inside your `UBER_HOME`
directory. You will need one file per host, app, host provider and app
provider. Apps can 'wrap' other apps, thus creating a 'stack of stacks' aka
an 'Uberstack'.

More information on configuring these files will be provided soon.

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
be expressed with a `docker-compose.yml` file, and enhanced with a 
`rancher-compose.yml` one. An "Uberstack" groups these stacks together,
and providing configuration (environment variables) for one or more 
'environments', such as local, dev, staging, qa and production.

Acknowledgements
----------------

 * Terraform configuration for AWS was inspired by this article:
   (https://www.airpair.com/aws/posts/ntiered-aws-docker-terraform-guide)
 * Inspiration for the VirtualBox support in Uberstack came from this article:
   (http://rancher.com/running-rancher-on-a-laptop/)

