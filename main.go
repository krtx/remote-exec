package main

import (
    "net/http"
    "log"
)

func main() {
    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        var err error
        switch r.Method {
        case "PUT":
            err = Put(w, r)
        case "POST":
            err = Post(w, r)
        case "DELETE":
            err = Delete(w, r)
        default:
            BadRequest(w)
        }
        if err != nil {
            log.Print(err)
            InternalServerError(w)
        }
    })
    log.Fatal(http.ListenAndServe(":8080", nil))
}
