**NOTE**: The project is in a proof of concept stage yet. Please do not use it in your production.

# What is Crankshaft?

Crankshaft is a barebone go web-server inspired by [unix philosophy](https://en.wikipedia.org/wiki/Unix_philosophy). 
It allows handling requests by chaining interceptors (sometimes referred ad middleware). 
These interceptors are golang [plugins](https://golang.org/pkg/plugin/) with one unified interface, 
that can be loaded and plugged into different routes. 
These plugged in interceptors along the request pipeline can then decide to handle the request, and 
break the chain or hand it over to next interceptor (may be modifying/adding stuff request/response).

## Philosophy

Dynamic modules is not a new concept and has been practiced for years in traditional web servers. In modern times Go 
due to it's efficiency and stability has been used by various projects to develop gateways/sidecars/web servers. Due
to single binary philosophy initially these projects provided no mechanism (except some scripting language) to extend
or provide interceptions, until recently. While these gateways are loaded with many pre-built interceptors/middleware 
to mold and modify request along the pipeline, crankshaft does the exact opposite. Crankshaft come with nothing builtin.
Everything has to be an interceptor, and interceptor in turn can be actually responding to request or chaining the 
request forward. 

## Why?

**In essence crankshaft is created to conceive unix philosophy interceptor plugins that can be chained**. This allows
composing complex request handling pipelines with simple, and clean chains.  

Crankshaft wants to take first principle approach on building a customizable gateway, writing minimal core will make 
it extremely testable, and stable. Due to intuitive nature and low friction writing interceptor plugin should be 
extremely easy. As an example; serving static file is not built-in to the server, a plugin `static_content`. Similarly
for path prefix stripping there is plugin `strip_path`. This should make it extremely simple to write plugins for
instrumentation, tracing, logging, authentication etc. 

# Show me some interceptor code

Interceptor in a plugin is a dead simple function that takes configuration `map[string]interface{}` and returns back
a function and error. Here is a very simple hello world interceptor:

```go
package main

import (
"fmt"
"log"
"net/http"
)

func Create(config map[string]interface{}) (func(http.ResponseWriter, *http.Request, http.HandlerFunc), error) {
    log.Println("Loading Hello World with", config, "...")
   
    return func(w http.ResponseWriter, r *http.Request, _ http.HandlerFunc) {
        fmt.Fprintf(w, "Greetings %s!", r.URL.Path[1:])
    }, nil
}

``` 

Compiling the plugin should be as simple as:

```
go build -buildmode=plugin
```

For loading `plugins.yml` will look something like this:
```yaml
hello_world: bin/modules/hello_world.so
```

Then define a chain invoking the plugin `chains.yaml`:
```yaml
- chain: hello_chain
  of:
    - plugin: hello_world
      config: {}
```

Finally plug this chain into router `routes.yaml`:
```yaml
- path: /
  invoke: hello_chain
```

Now running crankshaft `bin/crankshaft -chains chains.yaml -plugins plugins.yaml -routes routes.yaml` should
bring up a server, and hitting `http://localhost:6060/crankshaft` will return `Greetings crankshaft!`. 

# Yet to implement

 [ ] Github hooks
 [ ] Test case coverage
 [ ] Logging story
 [ ] Documentation 
 [ ] Zero downtime, hot reload
 [ ] TLS + HTTP2 support (partially present)
 [ ] Control pane
