# go-multishot

Simple HTTP repeater. Receives a GET/POST/other request and sends it to multiple servers.
Why use it?

* If you want to duplicate HTTP traffic from production system to a staging one
* If you need to test new production system with traffic from an existing one
* You can't use raw traffic duplication

## Usage

Compile with

    go build multishot.go

Then run

    ./multishot -downstreams SERVER1,SERVER2

multishot accepts a list of downstream servers to duplicate request to.
First downstream is treated as a 'main' server, and its response is returned from the multishot.

multishot accepts -port argument with port number to bind to.

## Monitoring

Special location /archer can be queried for:

* How many requests were handled
* How many cloned requests differed in status code from the origin

## Testing

I used endpoint.go for testing. It is a simple HTTP server that returns request content back.
Build and run two endpoints with

    go build endpoint.go
    ./endpoint -port ":8090" &
    ./endpoint -port ":8091" &
    ./multishot -downstreams localhost:8090,localhost:8091

Repeater will listen on the port 8080 and multiply each incoming request to both endpoints.

## TODO

* Improve logging
* Use httputil/ReverseProxy instead of hand-written version
* FIXME: On hign loads panic occurs when closing response body
* BUG: Improve error tolerance - do not crash when downstream is not active
* Improve /archer - monitor responses from different downstreams
** Monitor not only status code, but also SHA1 hash of the response
* Add config processing
** Downstreams list
** Comments
** Throttling ratio
