Uberstack Dev Stack Example
===========================

Please review the files in this directory. Note that each file
refers to either an app, a host, an app-provider or a host-provider. Each app
or host provider will require different configurations.

There are multiple of each type of provider. This example relies upon two:

 * TerraformHostProvider: this allows the creation of infrastructure from scratch,
   using Terraform to do so
 * DockerAppProvider: this provider will install apps specified by `docker-compose`
   files. It is intended to be used to bootstrap an infrastructure where all that
   exists on the client node is an installation of Docker.

The applications installed in this example are:

 * A Docker Registry - with SSL enabled via Letsencrypt and basic authentication
   enabled using credentials passed in as environment variables
 * A Rancher Server - also with basic authentication and logged into the Docker
   Registry
 * A Jenkins Server - again, with basic auth and logged into Docker Registry
 * An HTTP proxy to allow Rancher and Jenkins to both serve on port 80 but with
   unique domain names.

These will all be installed within an Amazon Web Services VPC that will be created
for you.

Partially complete code exists to install an OpenVPN VPN server within this VPC.

Usage
-----
To use this example, ensure the `uberstack` binary is on your path, and that the
`UBER_HOME` environment variable points at the `examples/dev-stack` directory.

This example will install the above apps to three domain names, docker.MYDOMAIN.com,
rancher.MYDOMAIN.com, and jenkins.MYDOMAIN.com. Before starting, inside AWS,
create an Elastic IP and take note of the "allocation ID". Point the above three
domain names at this new elastic IP. Wait for the domain names to resolve before
continuing.

Then, ensure that these environment variables are installed:

 * `UBER_DOMAIN`: the TLD to be used, e.g. example.com would give us
   docker.example.com, rancher.example.com and jenkins.example.com
 * `UBER_USER`: the username to use in all authentication
 * `UBER_PASS`: the password to use in all authentication
 * `UBER_EMAIL`: an email address to provide to Docker Registry
 * `EIP_ALLOCATION`: An Elastic IP allocation ID for an IP address created earlier
 * `PUBLIC_KEY`: The text of a public key that can be used to communicate with hosts

Start the host with:

    uberstack host up dev-host

Start the applications with:

    uberstack app up dev-stack dev

Note that 'dev-host' and 'dev-stack' here refer to the YAML files with the same name.

