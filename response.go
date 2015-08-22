package main

import (
    "net/http"
)

func MakeResponse(status int, message string, w http.ResponseWriter) {
    w.WriteHeader(status)
    w.Write([]byte(message))
}

func InternalServerError(w http.ResponseWriter) {
    MakeResponse(http.StatusInternalServerError, "Internal Server Error", w)
}

func BadRequest(w http.ResponseWriter) {
    MakeResponse(http.StatusBadRequest, "Bad Request", w)
}

func NotFound(w http.ResponseWriter) {
    MakeResponse(http.StatusNotFound, "Not found", w)
}
