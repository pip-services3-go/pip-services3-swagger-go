# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> Swagger UI for Golang

This module is a part of the [Pip.Services](http://pipservices.org) polyglot microservices toolkit.

The swagger module provides a Swagger UI that can be added into microservices and seamlessly integrated with existing REST and Commandable HTTP services.

The module contains the following packages:
- **Build** - Swagger service factory
- **Services** - Swagger UI service

<a name="links"></a> Quick links:

* [Change Log](CHANGELOG.md)
* [Get Help](https://www.pipservices.org/community/help)
* [Contribute](https://www.pipservices.org/community/contribute)


## Use

Install the Go package as
```bash
go get github.com/pip-services3-go/pip-services3-swagger-go
```

## Develop

For development you shall install the following prerequisites:
* Golang 1.13
* Visual Studio Code or another IDE of your choice
* Docker

Generate embedded resources:
```bash
./resgen.ps1
```

Run automated tests:
```bash
go test ./test/...
```

Generate API documentation:
```bash
./docgen.ps1
```

Before committing changes run dockerized build and test as:
```bash
./build.ps1
./test.ps1
./clear.ps1
```

## Contacts

The Golang version of Pip.Services is created and maintained by:
- **Sergey Seroukhov**
