package core

import (
    "log"
    "net/http"
    "plugin"
)

const CreateMethodName = "Create"

type LoadedPlugin func(map[string]interface{}) (func(http.ResponseWriter, *http.Request, http.HandlerFunc), error)

func loadPlugin(pluginName string) (LoadedPlugin, error) {
    log.Println("Loading plugin", pluginName, "...")
    plg, err := plugin.Open(pluginName)
    if err != nil {
        return nil, err
    }

    init, err := plg.Lookup(CreateMethodName)
    if err != nil {
        return nil, err
    }

    return init.(func(map[string]interface{}) (func(http.ResponseWriter, *http.Request, http.HandlerFunc), error)), nil
}

func LoadPlugins(mapping map[string]string) (map[string]LoadedPlugin, error) {
    ret := make(map[string]LoadedPlugin)
    for name, path := range mapping {
        plgn, err := loadPlugin(path)
        if err != nil {
            return nil, err
        }
        ret[name] = plgn
    }

    return ret, nil
}
