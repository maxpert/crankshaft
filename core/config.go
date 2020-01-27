package core

import (
    "io/ioutil"

    "gopkg.in/yaml.v2"
)

type ChainInterceptor struct {
    PluginName string                 `yaml:"plugin"`
    Config     map[string]interface{} `yaml:"config"`
}

type Chain struct {
    Name         string             `yaml:"chain"`
    Interceptors []ChainInterceptor `yaml:"of"`
}

type RouteSpecs struct {
    Host   string `yaml:"host"`
    Method string `yaml:"method"`
    Path   string `yaml:"path"`
    Invoke string `yaml:"invoke"`
}

type ChainsMap map[string][]ChainInterceptor
type LoadedPluginsMap map[string]LoadedPlugin
type Routes []RouteSpecs

func loadYaml(filePath string, unmarshalTo interface{}) error {
    configData, err := ioutil.ReadFile(filePath)
    if err != nil {
        return err
    }

    err = yaml.Unmarshal(configData, unmarshalTo)
    if err != nil {
        return err
    }

    return nil
}

func LoadRoutes(filePath string) (Routes, error) {
    if filePath == "" {
        return Routes{}, nil
    }

    ret := Routes{}
    if err := loadYaml(filePath, &ret); err != nil {
        return nil, err
    }

    return ret, nil
}

func LoadNamedChains(filePath string) (ChainsMap, error) {
    if filePath == "" {
        return ChainsMap{}, nil
    }

    chains := make([]Chain, 0)
    if err := loadYaml(filePath, &chains); err != nil {
        return nil, err
    }

    ret := ChainsMap{}
    for _, c := range chains {
        ret[c.Name] = c.Interceptors
    }

    return ret, nil
}

func LoadConfigPlugins(configPath string) (LoadedPluginsMap, error) {
    if configPath == "" {
        return LoadedPluginsMap{}, nil
    }

    mapping := make(map[string]string)
    if err := loadYaml(configPath, &mapping); err != nil {
        return nil, err
    }

    return LoadPlugins(mapping)
}
