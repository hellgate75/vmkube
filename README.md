<p align="center" style="width: 100%"><img width="200" height="200" src="/images/golang.png" />&nbsp;<img width="168" height="200" src="/images/docker-machine.png" /></p>

# Go Virtual Machine Kube (vmkube)

<p align="center"><img src="https://travis-ci.org/hellgate75/vmkube.svg?branch=master" alt="trevis-ci" width="98" height="20" />&nbsp;<a href="https://travis-ci.org/hellgate75/vmkube">Check last build on Travis-CI</a></p><br/>

Go-Lang Virtual Machine environment manager. This package allow to define projects and build infrastructures on local or remote cloud providers


## Prerequisites

* [Go](https://golang.org/dl/) (tested with version 1.12)
* [Docker](https://www.docker.com/get-docker) and [Docker-Machine](https://docs.docker.com/machine/install-machine/)

One of following :
* [VMWare Fusion](https://my.vmware.com/en/web/vmware/info/slug/desktop_end_user_computing/vmware_fusion/8_0)
* [VMWare Workstation Player](https://my.vmware.com/en/web/vmware/free#desktop_end_user_computing/vmware_workstation_player/12_0)
* [Virtual Box](https://www.virtualbox.org/wiki/Downloads) (VMWare Utilities on mac : `brew install vagrant-vmware-fusion`, on windows/linux: [VDDK65](https://my.vmware.com/group/vmware/get-download?downloadGroup=VDDK65))
* Microsoft [Hyper-V](https://www.manageengine.com/free-hyper-v-configuration/documents.html)
* Any of the cloud provider [supported by docker](https://docs.docker.com/machine/drivers/)

In order to define machine OSs, you can use any of the [supported ISOs](https://docs.docker.com/machine/drivers/os-base/)

## Goals

Define a virtual machine remote/local manager based on docker-machine drivers. It defines a smart and simple command-line/web interface to manage environments.
It provides features to define project and deployment plans for infrastructures and applications.
Domain, Network and Host management level and hierarchy are main concepts in the VMKube philosophy.
WMKube provides development and deployment environments: You have to define an initial project, you can manager the networks, domains, hosts and a staging phase.
When you close the project, you are ready to delete it or, alternatively, to build and run the infrastructure.

## What is provided?

Provided features:

* Project definition procedures

* Project staging/un-staging procedures

* Project build procedures

* Infrastructure build procedures

* Infrastructure/Project inspection procedures

* Custom Deployment plans with main providers (VMKubelet, [Ansible](http://docs.ansible.com/ansible/index.html), [Helm for Kubernetes](https://github.com/kubernetes/helm), ...)

* Digital Control and multi vendor instances (Machine and Cloud-Machines can be defined)

* Multiple project/infrastructure information export formats

Machine Providers:

See [Docker-Machine Drivers](https://docs.docker.com/machine/drivers/)

To compile and run this project you have to check availability of following software:
* [Go](https://golang.org/dl/) (tested with version 1.8)
* [Docker](https://www.docker.com/get-docker) and [Docker-Machine](https://docs.docker.com/machine/install-machine/)
* Test and Utility GOLang packages ([UUID Package](https://github.com/satori/go.uuid), [Unit Test](https://github.com/stretchr/testify)) and [GO SSH Terminal](http://golang.org/x/crypto/ssh/terminal), [YAML Parser](http://gopkg.in/yaml.v2)


## Architecture



## Configuration


## Checkout and test this repository

Go in you `GOPATH/src` folder and type :
```sh
 go get github.com/stretchr/testify
 go get github.com/satori/go.uuid
 go get golang.org/x/crypto/ssh/terminal
 go get gopkg.in/yaml.v2
 git clone https://github.com/hellgate75/vmkube.git

```
or simply :
```sh
 go get github.com/stretchr/testify
 go get github.com/satori/go.uuid
 go get golang.org/x/crypto/ssh/terminal
 go get gopkg.in/yaml.v2
 go get github.com/hellgate75/vmkube
```


## Build

It's present a make file that returns an help on the call :

```sh
make
```
Provided `Makefile` help returns following options :
```sh
make [all|init|test|build|exe|run|clean|install]
all: test build exe run
init: get required external packages
test: run unit test
build: build the module
exe: make executable for the module
clean: clean module C objects
run: exec the module code
install: install the module in go libs
```

Alternatively you can execute following commands :
 * `go get github.com/stretchr/testify` to download unit test external package
 * `go get github.com/satori/go.uuid` to download UUID management external package
 * `go get golang.org/x/crypto/ssh/terminal` to download SSH terminal external package
 * `go get gopkg.in/yaml.v2` to download YAML parser
 * `go build .` to build the project
 * `go test` to run unit and integration test on the project
 * `go run main.go` to execute the project
 * `go build --buildmode exe .` to create an executable command
 * `go install` to install the executable command


## Execution

The tool provides an help section, describing commands, sub-commands and has a nested help level for commands details.

The help is available executing :
* `vmkube help` General Help
* `vmkube help [command]` Detailed Command syntax helper

Import / Alter Project Commands provides a sample for the expected input format. Import and Export of components is provided in following file formats:
* JSON - standard JSON language
* XML - Untagged and un-described XML format (Pure XML tag sequence, no XML definition or version is accepted).
* YAML - standard YAML format.

In this release the command list is composed by following keys :
* `help` : Show generic commands help
* `start-infra` : Start an existing Infrastructure
* `stop-infra` : Stop a Running Infrastructure
* `restart-infra` : Restart a Running Infrastructure
* `destroy-infra` : Destroy a specific Infrastructure
* `backup-infra` : Backup a specific Infrastructure to a backup file
* `recover-infra` : Recover a specific Infrastructure from a backup file
* `infra-status` : Require information about a specific Infrastructure
* `list-all-infra` : Require list of all Infrastructures
* `list-projects` : Require list of all available projects
* `project-status` : Require information about a specific projects
* `define-project` : Creates a new project
* `alter-project` : Modify existing project, e.g.: open, close project or add, modify, delete items
* `info-project` : Provides information about project elements definition
* `delete-project` : Delete existing project
* `build-project` : Build and existing project and create/modify an infrastructure
* `import-project` : Import project from existing configuration
* `export-project` : Export existing project configuration


## Author
[Fabrizio Torelli](https://ie.linkedin.com/in/fabriziotorelli) is Cloud/System Architect working in the IT sector since 1999.

## License

Licensed under the [MIT](/LICENSE) License (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[https://opensource.org/licenses/MIT](https://opensource.org/licenses/MIT)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
