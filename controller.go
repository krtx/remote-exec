package main

import (
    "net/http"
    "os"
    "encoding/json"
)

func Post(w http.ResponseWriter, r *http.Request) error {
    r.ParseMultipartForm(32 << 20)
    file, _, err := r.FormFile("source")
    if err == http.ErrNotMultipart {
        BadRequest(w)
        return nil
    }
    if err != nil {
        return err
    }
    defer file.Close()
    id, err := Save(file)
    if err != nil {
        return err
    }
    w.Header().Set("Content-type", "text/plain")
    w.Write([]byte(id))
    return nil
}

func Put(w http.ResponseWriter, r *http.Request) error {
    id := r.FormValue("id")
    if !Validate(id) {
        BadRequest(w)
        return nil
    }
    if !Exist(id) {
        NotFound(w)
        return nil
    }
    cmd := r.FormValue("command")
    args := r.FormValue("args")
    sout, serr, status, err := RunCmd(id, cmd, args)
    if err != nil {
        return err
    }
    b, err := json.Marshal(map[string]string{
        "stdout": sout,
        "stderr": serr,
        "status": status,
    })
    if err != nil {
        return err
    }
    w.Header().Set("Content-type", "text/json")
    w.Write(b)
    return nil
}

func Delete(w http.ResponseWriter, r *http.Request) error {
    id := r.FormValue("id")
    if !Validate(id) {
        BadRequest(w)
        return nil
    }
    if !Exist(id) {
        NotFound(w)
        return nil
    }
    return os.RemoveAll(BaseDir + "/" + id)
}
