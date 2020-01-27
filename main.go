package main

import (
    "context"
    "flag"
    "log"
    "net/http"

    "crankshaft/core"
    "golang.org/x/net/http2"
    "golang.org/x/net/http2/h2c"
)

var (
    pluginsConfigPath = flag.String(
        "plugins",
        "",
        "Config yaml defining the plugin mapping",
    )

    chainsConfigPath = flag.String(
        "chains",
        "",
        "Chains yaml defining all middleware chains",
    )

    routesConfigPath = flag.String(
        "routes",
        "",
        "Routes for configuration server",
    )

    bindAddress = flag.String(
        "bind",
        ":6060",
        "Bind address of the server",
    )
)

func buildServerHandler(p core.LoadedPluginsMap, c core.ChainsMap, r core.Routes) http.HandlerFunc {
    chains, err := core.BuildChains(p, c)

    if err != nil {
        log.Panicln(err)
    }

    log.Println("Linking routes", r, "to chains", chains, "...")
    routes := make([]*core.RouteInvoker, 0)
    for _, routeSpecs := range r {
        ch, ok := chains[routeSpecs.Invoke]
        if !ok {
            log.Panicf("Unable to find chain %s for route %s", routeSpecs.Invoke, routeSpecs.Path)
        }

        routes = append(routes, core.NewRouteInvoker(routeSpecs, ch))
    }

    return func(w http.ResponseWriter, r *http.Request) {

        // Shared context metadata for whole middleware chain
        sharedMetaData := make(map[interface{}]interface{})
        r = r.WithContext(context.WithValue(r.Context(), "shared-meta", sharedMetaData))

        for _, i := range routes {
            if i.Invoke(w, r) {
                return
            }
        }

        http.Error(w, "Unable to match any route", 404)
    }
}

func main() {
    flag.Parse()

    chains, err := core.LoadNamedChains(*chainsConfigPath)
    if err != nil {
        log.Panicln("Unable to load chains", err)
    }

    routes, err := core.LoadRoutes(*routesConfigPath)
    if err != nil {
        log.Panicln("Unable to load routes", err)
    }

    plugins, err := core.LoadConfigPlugins(*pluginsConfigPath)
    if err != nil {
        log.Panicln("Unable to load plugins", err)
    }

    handlerFunc := buildServerHandler(plugins, chains, routes)

    // Create and configure http server
    http2Server := &http2.Server{}
    server := &http.Server{}
    if err := http2.ConfigureServer(server, http2Server); err != nil {
        log.Panicln(err)
    }

    server.Handler = h2c.NewHandler(handlerFunc, http2Server)
    server.Addr = *bindAddress

    log.Println("Starting server...")
    err = server.ListenAndServe()
    if err != nil {
        log.Panicln(err)
    }
}
