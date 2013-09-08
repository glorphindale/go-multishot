# go-multishot

Simple http repeater. Handles request and sends it to multiple servers.

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

* Improve logging
* Use httputil/ReverseProxy instead of hand-written version

