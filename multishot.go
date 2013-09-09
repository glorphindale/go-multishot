package main

import (
    "bytes"
    "flag"
    "fmt"
    "io"
    "log"
    "net/http"
    "strings"
)

var downstreams []string

func forward_request(downstream string, in_req http.Request, body io.Reader) (resp *http.Response, err error) {
    client := &http.Client{}

    // Clone incoming request
    out_url := "http://" + downstream + in_req.URL.Path
    req, _ := http.NewRequest(in_req.Method, out_url, body)
    for h, vv := range in_req.Header {
        for _, v := range vv {
            req.Header.Add(h, v)
        }
    }
    req.ContentLength = in_req.ContentLength

    resp, err = client.Do(req) // TODO need to close r.Body, see http://golang.org/pkg/net/http/#Client.Do
    if err != nil {
        log.Println("Request", req.URL.Path, " to downstream", downstream, "failed", err)
        return
    }
    return
}

func handler(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request for", r.URL.Path)

    // Read the body, make it available for all the downstreams
    clength := r.ContentLength
    var raw_body []byte

    if clength != 0 {
        body_part := make([]byte, clength)
        n, err := r.Body.Read(body_part)
        if err == nil {
            raw_body = body_part[:n]
        }
    }

    for _, downstream := range downstreams[1:] {
        log.Println("Firing off", downstream)
        go forward_request(downstream, *r, bytes.NewReader(raw_body))
    }

    log.Println("Firing main", downstreams[0])
    resp, err := forward_request(downstreams[0], *r, bytes.NewReader(raw_body))
    if err == nil {
        io.Copy(w, resp.Body)
    } else {
        log.Println("Error", err)
        w.WriteHeader(503)
    }
}

func main() {
    var port = flag.Int("port", 8080, "port to listen on")
    var downstreams_raw = flag.String("downstreams", "localhost:8090,localhost:8091", "list of downstreams in 'host:port' format, separated by comma")
    flag.Parse()

    downstreams = strings.Split(*downstreams_raw, ",")

    http.HandleFunc("/", handler)
    port_string := fmt.Sprintf(":%d", *port)
    log.Println("Listening on port", port_string)
    http.ListenAndServe(port_string, nil)
}

