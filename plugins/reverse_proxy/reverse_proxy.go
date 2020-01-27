package main

import (
    "errors"
    "fmt"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "strings"
    "sync/atomic"
    "time"
)

type reverseProxy struct {
    urls    []*url.URL
    proxies []*httputil.ReverseProxy
    index   uint32
}

func loadUrls(config map[string]interface{}) ([]*url.URL, error) {
    hostsConfig, ok := config["urls"]
    if !ok {
        return nil, errors.New("please specify reverse proxy hosts")
    }

    hostsString, ok := hostsConfig.(string)
    if !ok {
        return nil, fmt.Errorf("hosts should be a string found: %v", hostsConfig)
    }

    hosts := strings.Split(hostsString, ", ")
    for i, h := range hosts {
        hosts[i] = strings.TrimSpace(h)
    }

    ret := make([]*url.URL, 0)
    for _, h := range hosts {
        parsedUrl, err := url.Parse(h)
        if err != nil {
            return nil, err
        }

        ret = append(ret, parsedUrl)
    }

    return ret, nil
}

func newReverseProxy(urls []*url.URL) *reverseProxy {
    proxies := make([]*httputil.ReverseProxy, 0)
    for _, u := range urls {
        p := httputil.NewSingleHostReverseProxy(u)
        proxies = append(proxies, p)
    }
    ret := &reverseProxy{
        urls:    urls,
        proxies: proxies,
    }

    for index := range ret.proxies {
        ret.proxies[index].Director = ret.requestDirector(index, ret.proxies[index].Director)
        ret.proxies[index].ModifyResponse = ret.modifyResponse(index, ret.proxies[index].ModifyResponse)
    }

    return ret
}

func (r *reverseProxy) modifyResponse(_ int, parent func(r *http.Response) error) func(r *http.Response) error {
    return func(r *http.Response) error {
        if parent != nil {
            if err := parent(r); err != nil {
                return err
            }
        }

        r.Header.Add("X-Handled-By", r.Request.URL.String())
        r.Header.Add("X-Enter-Time", r.Request.Header.Get("X-Enter-Time"))
        r.Header.Add("X-Exit-Time", time.Now().Format(time.RFC3339Nano))
        return nil
    }
}

func (r *reverseProxy) requestDirector(index int, parent func(*http.Request)) func(r *http.Request) {
    u := r.urls[index]
    return func(r *http.Request) {
        if parent != nil {
            parent(r)
        }

        r.Header.Add("X-Forwarded-Host", u.Host)
        r.Header.Add("X-Origin-Host", r.Host)
        r.Header.Add("X-Enter-Time", time.Now().Format(time.RFC3339Nano))
        r.Host = u.Host
    }
}

func (r *reverseProxy) handle(rw http.ResponseWriter, req *http.Request, _ http.HandlerFunc) {
    index := atomic.AddUint32(&r.index, 1)
    index = index % uint32(len(r.proxies))
    p := r.proxies[index]
    p.ServeHTTP(rw, req)
}

func Create(config map[string]interface{}) (func(http.ResponseWriter, *http.Request, http.HandlerFunc), error) {
    log.Println("Creating reverse proxy with config", config, "...")
    urls, err := loadUrls(config)
    if err != nil {
        return nil, err
    }

    rp := newReverseProxy(urls)
    return rp.handle, nil
}
