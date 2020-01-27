package core

import (
    "errors"
    "fmt"

    "crankshaft/api"
)

func BuildChains(p LoadedPluginsMap, c ChainsMap) (map[string]api.InterceptorChain, error) {
    ret := make(map[string]api.InterceptorChain)
    for name, interceptors := range c {
        chain := api.InterceptorChain{}
        for _, interceptor := range interceptors {
            if installer, ok := p[interceptor.PluginName]; !ok {
                return nil, errors.New(
                    fmt.Sprintf("Unknown plugin %s request for chain %s", interceptor.PluginName, name),
                )
            } else if plg, err := installer(interceptor.Config); err != nil {
                return nil, err
            } else {
                chain = append(chain, plg)
            }
        }

        ret[name] = chain
    }

    return ret, nil
}
