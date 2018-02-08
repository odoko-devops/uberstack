UberStack
=========

UberStack is a tool for describe applications by detailing their component parts.

With the rise of containers, deploying complex micro-service architectures has become
a lot simpler than it once was. However, we can still find ourselves thinking a lot
about containers and services, and not enough about applications.

UberStack's purpose is to allow us to specify a whole application by combining
service definitions.

In this current implementation, it makes use of Rancher, and thus applications
are made of docker-compose.yml and rancher-compose.yml files.

Why "UberStack"? It is common to have stacks as a grouping of services in containerised
technologies. An uberstack is a stack of stacks - allowing for multiple levels of
composition. For example, we may have a "mongo" stack that we then include within
various applications.

## Version Compatibility
UberStack requires Rancher CLI and is compatible with Rancher 1.6. 

## Running

You can check the [releases page](https://github.com/odoko-devops/uberstack/releases) for direct downloads of the binary or [build your own](#building). 

## Setting 

Further details on setting up UberStack will be forthcoming.

## Building

The binaries will be located in `/bin`.

### Linux binary

Run `make`.

### Mac binary

Run `make mac`

## Contact

For bugs, questions, comments, corrections, suggestions, etc., open an issue in
[odoko-devops/uberstack](//github.com/odoko-devops/uberstack/issues).

Or just [click here](//github.com/rancher/rancher/issues/new?title=%5Bcli%5D%20) to create a new issue.

## License
Copyright (c) 2018 [Odoko Ltd](http://www.odoko.com)

Contains code that is (c) 2014-2016 [Rancher Labs, Inc](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
