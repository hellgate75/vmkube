<p align="center" style="width: 100%"><img width="200" height="200" src="/images/golang.png" />&nbsp;<img width="168" height="200" src="/images/docker-machine.png" /></p>

# Go Virtual Machine Kube (vmkube)

Go-Lang Virtual environment manager. This package allow to define projects and build infrastructures on local or remote environments


## Goals

Define a virtual machine remote/local manager based on docker-machine drives.


## Pre-requisites

To compile and run this project you have to check availability of following software:
* [Go](https://golang.org/dl/) (tested with version 1.8)
* Docker and Docker-Machine
* Test and Utility GOLang packages (github.com/satori/go.uuid and github.com/stretchr/testify)


## Architecture



## Configuration


## Checkout and test this repository

Go in you `GOPATH\src` folder and type :
```sh
 git clone https://github.com/hellgate75/vmkube.git

```
or simply :
```sh
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
 * `go build .` to build the project
 * `go test` to run unit and integration test on the project
 * `go run main.go` to execute the project
 * `go build --buildmode exe .` to create an executable command


## Further test 




## Execution



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
