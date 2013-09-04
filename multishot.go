package main

import (
    "io"
    "log"
    "net/http"
    "strings"
)

var downstreams = []string{"localhost:8090", "localhost:8091"}

func forward_request(downstream string, r *http.Request) (resp *http.Response, err error) {
    client := &http.Client{}

    out_url := "http://" + downstream + r.URL.Path
    req, err := http.NewRequest("GET", out_url, r.Body) // Need to handle POST reqiests
    if err != nil {
        log.Println("Downstream", downstream, "failed")
        return
    }

    for h, v := range r.Header {
        req.Header.Add(h, strings.Join(v, ";"))
    }

    resp, err = client.Do(req)
    return
}

func handler(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request for", r.URL.Path)
    for _, downstream := range downstreams[1:] {
        log.Println("Firing off", downstream)
        go forward_request(downstream, r)
    }

    log.Println("Firing main", downstreams[0])
    resp, _ := forward_request(downstreams[0], r)
    io.Copy(w, resp.Body)
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}

