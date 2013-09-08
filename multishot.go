package main

import (
    "flag"
    "fmt"
    "io"
    "log"
    "net/http"
    "strings"
)

var downstreams []string

func forward_request(downstream string, r *http.Request) (resp *http.Response, err error) {
    client := &http.Client{}

    out_url := "http://" + downstream + r.URL.Path
    req, _ := http.NewRequest("GET", out_url, r.Body) // Need to handle POST reqiests
    for h, v := range r.Header {
        req.Header.Add(h, strings.Join(v, ";"))
    }
    resp, err = client.Do(req)
    if err != nil {
        log.Println("Request", r.URL.Path, " to downstream", downstream, "failed")
        return
    }
    return
}

func handler(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request for", r.URL.Path)
    for _, downstream := range downstreams[1:] {
        log.Println("Firing off", downstream)
        go forward_request(downstream, r)
    }

    log.Println("Firing main", downstreams[0])
    resp, err := forward_request(downstreams[0], r)
    if err == nil {
        io.Copy(w, resp.Body)
    } else {
        w.WriteHeader(resp.StatusCode)
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

