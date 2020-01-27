package main

import (
    "log"
    "net/http"
)

func Create(config map[string]interface{}) (func(http.ResponseWriter, *http.Request, http.HandlerFunc), error) {
    log.Println("Creating strip path plugin with config", config, "...")
    stripPath := "/"
    if p, ok := config["prefix"]; ok {
        stripPath = p.(string)
    }

    return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
        http.StripPrefix(stripPath, next).ServeHTTP(w, r)
    }, nil
}