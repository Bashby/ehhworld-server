## EhhWorld Server

<p align="center">
  <img src="assets\image\example_generated_v3.png" width="800" height="800" title="Mazeclient Logo" style="border-radius: 5px;">
</p>

The server logic for EhhWorld.

- Written in [golang](https://golang.org/).
- Uses [dep](https://github.com/golang/dep) for dependency management.

# Getting started

Before building:

- Install [golang](https://golang.org/)
- Install [dep](https://github.com/golang/dep)

Then, from the root of the project:

    cd server
    dep ensure

Finally, use one of the following build methods.

## Local Building

From the root of the project:

    cd server
    go run main.go --mode 1 --serve

## Docker Building

From the root of the project:

    cd deploy
    ./start-services.sh

## Interacting with the server

In your browser, navigate to `localhost:8081`.

## Notes

You can run `protoc --go_out=. *.proto` to generate the protocol buffers.
