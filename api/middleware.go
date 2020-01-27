package api

import "net/http"

type Interceptor func(http.ResponseWriter, *http.Request, http.HandlerFunc)

type InterceptorFunc http.HandlerFunc

type InterceptorChain []Interceptor

func (cont InterceptorFunc) Intercept(mw Interceptor) InterceptorFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        mw(writer, request, http.HandlerFunc(cont))
    }
}

func (chain InterceptorChain) Handler(handler http.HandlerFunc) http.Handler {
    curr := InterceptorFunc(handler)
    for i := len(chain) - 1; i >= 0; i-- {
        mw := chain[i]
        curr = curr.Intercept(mw)
    }

    return http.HandlerFunc(curr)
}
