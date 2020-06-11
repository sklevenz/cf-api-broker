# Cloud Foundry API Broker

This is an [OSB Service Broker](https://www.openservicebrokerapi.org/) for the [Cloud Foundry](https://www.cloudfoundry.org/) API endpoint.  

## Build and Development Process

This project uses a Python build script. For help call ´./bin/make.sh help´.

|          Make          | Description                                              |
|:----------------------:|----------------------------------------------------------|
| ./bin/make.sh test     | Call go test and go vet for all packages                 |
| ./bin/make.sh build    | Call go clean, go fmt, go build and set version          |
| ./bin/make.sh run      | Call go run to start server on port 5000 and set version |
| ./bin/make.sh generate | Generate OSB model API (we do not use the server)        |
| ./bin/make.sh release  | Create a release with goreleaser                         |



### Build Tools Required

- go programming language
- wget
- openapi-generate
- goreleaser
  
For OSX: `brew install wget openapi-generator goreleaser`

