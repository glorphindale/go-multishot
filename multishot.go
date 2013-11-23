package main

import (
    "bytes"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "net/http"
    "strings"
    "sync/atomic"
)

var downstreams []string
var request_id uint64
var mismatched_count uint64

type DownstreamResult struct {
    name string
    code int
}

// Forward *in_req* to specified *downstream*, with *body* (when POST is handled)
// if *close_body* is true - close the body before returning
func forward_request(downstream string, in_req http.Request, body io.Reader, close_body bool, comm chan DownstreamResult) (resp *http.Response, err error) {
    client := &http.Client{}

    // Clone incoming request
    out_url := "http://" + downstream + in_req.URL.Path
    if in_req.URL.RawQuery != "" {
        out_url += "?" + in_req.URL.RawQuery;
    }
    req, _ := http.NewRequest(in_req.Method, out_url, body)
    for h, vv := range in_req.Header {
        for _, v := range vv {
            req.Header.Add(h, v)
        }
    }
    req.ContentLength = in_req.ContentLength

    resp, err = client.Do(req) // TODO need to close r.Body, see http://golang.org/pkg/net/http/#Client.Do
    comm <- DownstreamResult{downstream, resp.StatusCode}
    if err != nil {
        log.Fatal("Request ", req.URL.Path, " to downstream ", downstream, " failed ", err)
        return
    }
    if close_body {
        resp.Body.Close()
    }
    return
}

func stat_handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Handled %d requests\n", request_id)
    fmt.Fprintf(w, "Mismatched %d responses\n", mismatched_count)
}

func exit_handler(w http.ResponseWriter, r *http.Request) {
    os.Exit(0)
}

func collector(count int, in chan DownstreamResult) {
    storage := make([]DownstreamResult, count)

    // This assumes that main downstream responds first
    // TODO properly detect main downstream
    storage[0] = <-in
    fmt.Println("Main downstream respondes with ", storage[0].code)

    for i := 1; i < count; i++ {
        storage[i] = <-in
        fmt.Println("Main downstream respondes with ", storage[0].code)
        if storage[i].code != storage[0].code {
            fmt.Println("Downstream ", storage[i].name, " mismatch:", storage[i].code)
            atomic.AddUint64(&mismatched_count, 1)
        }
    }
}

func main_handler(w http.ResponseWriter, r *http.Request) {
    rid := atomic.AddUint64(&request_id, 1)
    log.Println("Received request", rid, r.URL.Path, r.URL.RawQuery)

    comm := make(chan DownstreamResult)
    go collector(len(downstreams), comm)

    // Read the body, make it available for all the downstreams
    clength := r.ContentLength
    var raw_body []byte

    if clength != 0 {
        body_part := make([]byte, clength)
        n, err := r.Body.Read(body_part)
        if err == nil {
            raw_body = body_part[:n]
        }
        r.Body.Close()
    }

    for _, downstream := range downstreams[1:] {
        go forward_request(downstream, *r, bytes.NewReader(raw_body), true, comm)
    }

    resp, err := forward_request(downstreams[0], *r, bytes.NewReader(raw_body), false, comm)
    if err == nil {
        io.Copy(w, resp.Body)
        resp.Body.Close()
    } else {
        w.WriteHeader(503)
    }
}

func main() {
    var port = flag.Int("port", 8080, "port to listen on")
    var downstreams_raw = flag.String("downstreams", "localhost:8090,localhost:8091", "list of downstreams in 'host:port' format, separated by comma")
    flag.Parse()

    downstreams = strings.Split(*downstreams_raw, ",")

    http.HandleFunc("/archer", stat_handler)
    http.HandleFunc("/exit", exit_handler)
    http.HandleFunc("/", main_handler)
    port_string := fmt.Sprintf(":%d", *port)
    log.Println("Listening on port", port_string)
    http.ListenAndServe(port_string, nil)
}

