package main

import (
    "log"
    "net/http"
)

var server http.Handler = nil

func Create(config map[string]interface{}) (func(http.ResponseWriter, *http.Request, http.HandlerFunc), error) {
    log.Println("Creating static content plugin with config", config, "...")
    directoryPath := "static"
    if p, ok := config["root"]; ok {
        directoryPath = p.(string)
    }

    server = http.FileServer(http.Dir(directoryPath))
    return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
        server.ServeHTTP(w, r)
    }, nil
}