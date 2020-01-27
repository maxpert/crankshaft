package core

import (
    "net/http"
    "strings"

    "crankshaft/api"
)

type RouteInvoker struct {
    specs RouteSpecs
    chain api.InterceptorChain
}

func NewRouteInvoker(specs RouteSpecs, chain api.InterceptorChain) *RouteInvoker {
    return &RouteInvoker{
        specs: specs,
        chain: chain,
    }
}

func (i *RouteInvoker) isHostMatch(r *http.Request) bool {
    if i.specs.Host == "" || i.specs.Host == "*" {
        return true
    }

    return i.specs.Host == r.Host
}

func (i *RouteInvoker) isMethodMatch(r *http.Request) bool {
    if i.specs.Method == "" || i.specs.Method == "*" {
        return true
    }

    return i.specs.Method == r.Method
}

func (i *RouteInvoker) isPathMatch(r *http.Request) bool {
    if i.specs.Path == "" || i.specs.Path == "/" {
        return true
    }

    return strings.HasPrefix(r.URL.Path, i.specs.Path)
}

func (i *RouteInvoker) handler(w http.ResponseWriter, r *http.Request) {
    http.Error(w, "No terminating middleware", 404)
}

func (i *RouteInvoker) invokeChain(w http.ResponseWriter, r *http.Request) {
    i.chain.Handler(i.handler).ServeHTTP(w, r)
}

func (i *RouteInvoker) Invoke(w http.ResponseWriter, r *http.Request) bool {
    if i.isHostMatch(r) && i.isMethodMatch(r) && i.isPathMatch(r) {
        i.invokeChain(w, r)
        return true
    }

    return false
}
