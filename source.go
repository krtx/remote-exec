package main

import (
    "time"
    "io"
    "os"
    "path/filepath"
    "math/rand"
    "archive/tar"
    "regexp"
    "os/exec"
    "bytes"
    "strings"
)

const (
    BaseDir = "d"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func RandomString(n int) string {
    rand.Seed(time.Now().UnixNano())
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func Save(ir io.Reader) (string, error) {
    id := RandomString(32)
    tr := tar.NewReader(ir)
    dir := BaseDir + "/" + id + "/"
    for {
        hdr, err := tr.Next()
        if err == io.EOF {
            break
        }
        if err != nil {
            return "", err
        }
        path := dir + hdr.Name
        if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
            return "", err
        }
        if path[len(path) - 1] == '/' {
            continue
        }
        f, err := os.Create(path)
        if err != nil {
            return "", err
        }
        defer f.Close()
        if _, err := io.Copy(f, tr); err != nil {
            return "", err
        }
    }
    return id, nil
}

func Validate(id string) bool {
    if matched, _ := regexp.Match("[^a-z0-9]", []byte(id)); matched {
        return false
    }
    return true
}

func Exist(id string) bool {
    _, err := os.Stat(BaseDir + "/" + id)
    return err == nil
}

func RunCmd(id string, cmd string, argsStr string) (string, string, string, error) {
    proc := exec.Command(cmd)
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    proc.Stdout = &stdout
    proc.Stderr = &stderr
    proc.Dir = BaseDir + "/" + id
    if len(argsStr) > 0 {
        args := strings.Split(argsStr, " ")
        for i := range args {
            proc.Args = append(proc.Args, args[i])
        }
    }
    err := proc.Run()
    if err != nil {
        if _, ok := err.(*exec.ExitError); !ok {
            return "", "", "", err
        }
    }
    return stdout.String(), stderr.String(), proc.ProcessState.String(), nil
}
