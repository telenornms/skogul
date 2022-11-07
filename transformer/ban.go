package transformer

import (
    "errors"
    "github.com/telenornms/skogul"
    "strings"
)

type Ban struct {
    DPaths map[string]interface{} `doc:"Data map of key value pairs where the keys are . separated tree paths e.g foo.bar.baz: true"`
    MPaths map[string]interface{} `doc:"Metadata map of key value pairs where the keys are . separated tree paths e.g foo.bar.baz: true"`
}

func (b *Ban) Transform(c *skogul.Container) error {
    var err error

    for _, mi := range c.Metrics {
        for pathKey, pathValue := range b.DPaths {
            splittedPath := strings.Split(pathKey, ".")
            if _, ok := mi.Data[splittedPath[0]]; ok {
                mi.Data, err = b.traverseDepths(mi.Data, splittedPath, pathValue, 0)
            }
        }

        for pathKey, pathValue := range b.MPaths {
            splittedPath := strings.Split(pathKey, ".")
            if _, ok := mi.Metadata[splittedPath[0]]; ok {
                mi.Metadata, err = b.traverseDepths(mi.Metadata, splittedPath, pathValue, 0)
            }
        }
    }

    return err
}

/*
    Recursively traverse a nested tree of elements based on path and remove last element in the tree
*/
func (b *Ban) traverseDepths(d map[string]interface{}, path []string, pathValue interface{}, depth int) (map[string]interface{}, error) {
    var err error
    if depth == len(path) - 1 {
        if d[path[len(path) - 1]] == pathValue {
            delete(d, path[len(path) - 1])
            return d, err
        }
        return d, err
    }

    if _, ok := d[path[depth]]; ok {
        key, okk := d[path[depth]].(map[string]interface{})
        if !okk {
            return d, errors.New("could not cast key to map")
        }

        d[path[depth]], err = b.traverseDepths(key, path, pathValue, depth + 1)
        return d, err
    }

    return d, errors.New("invalid key occurred")
}