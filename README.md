# go-multishot
============

Simple http multishot. Handles request and sends it to multiple servers.

## Usage

Compile with
    go build endpoint.go
    go build multishot.go

Run two endpoints with
    ./endpoint --port ":8090" &
    ./endpoint --port ":8091" &

Then run
    ./multishot

Repeater will listen on the port 8080 and multiply each incoming GET request to both endpoints.

## TODO

* Add config file for multishot
* Handle POST requests
* Improve logging

