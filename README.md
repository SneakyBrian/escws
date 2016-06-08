# ESCWS

## Embedded Static Content Web Server

This simple application embeds the contents of the /static/ folder into the executable binary and then serves the files using the applications built in web server.

## Getting

Get the source:

`go get github.com/sneakybrian/escws`

## Building

From within the source directory:

`go generate`

Embeds the contents of the ./static directory into the executable binary

`go build`

Builds the executable binary

And optionally:

`go install`

To install to your $GOPATH$\bin directory

## Development

The program has a test page available at:

http://localhost:8181/static/test.html

The program also has a configurable port number for the HTTP Server:

`escws -port=1234`

Would run the web server at:

http://localhost:1234/

Run the program with the `-useLocal` flag in order to serve the static resources from the file system rather than the embedded resources:

`escws -useLocal`

This allows you to edit the contents of the ./static folder on-the-fly and be able to refresh the browser to get the updated changes

## Packages

Uses the excellent `esc` package for generating the embedded static resources:

https://github.com/mjibson/esc