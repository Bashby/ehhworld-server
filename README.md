# README #

This is the logic to generate maps for EHH.IO.

Written in [golang](https://golang.org/).
Uses [dep](https://github.com/golang/dep) for dependency management.

## Getting started

Before building:
    - Install [golang](https://golang.org/)
    - Install [dep](https://github.com/golang/dep)

From the root of the project:
    - `cd server`
    - `dep ensure`

Then, do one of the following build methods.

### Local Building

From the root of the project:
    - `cd server`
    - `go run main.go --mode 1 --serve`

### Docker Building

From the root of the project:
    `cd deploy`
    `./start-services.sh`

## Interacting with the server

Navigate to `localhost:8080`.

# Dev Notes

You can run `protoc --go_out=. *.proto` to generate the protocol buffers.