package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
    msg := "Request " + r.URL.Path
    for k, v := range r.Header {
        msg += ";\n\t" + k + ": " + strings.Join(v, ",")
    }

    // Copy the request body
    body_part := make([]byte, 20)
    var body string
    n, err := r.Body.Read(body_part)
    if err != nil {
        body = "<Empty>"
    } else {
        body = string(body_part[:n])
    }

    msg += "\nBody: '" + body + "'"

    log.Println(msg)
    fmt.Fprintf(w, msg)
}

func main() {
    var port = flag.String("port", ":8090", "Port to listen on")
    flag.Parse()

    log.Println("Listening on", *port)

    http.HandleFunc("/", handler)
    http.ListenAndServe(*port, nil)
}

