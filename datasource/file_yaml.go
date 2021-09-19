package datasource

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type YamlFile map[string]string

func NewYamlFile(file string) (YamlFile, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	y := make(YamlFile)
	for k, v := range m {
		flatten(y, k, v)
	}

	datasources = append(datasources, y)

	return y, nil
}

func flatten(result YamlFile, key string, val interface{}) {
	switch value := val.(type) {
	case map[interface{}]interface{}:
		for k, v := range value {
			flatten(result, fmt.Sprintf("%s.%v", key, k), v)
		}
		return
	case []interface{}:
		for i, v := range value {
			flatten(result, fmt.Sprintf("%s.%d", key, i), v)
		}
		return
	}
	result[key] = fmt.Sprintf("%v", val)
}

func (y YamlFile) Get(key string) (string, bool) {
	v, ok := y[key]
	return v, ok
}
