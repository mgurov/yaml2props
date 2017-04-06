package pkg

import (
	"bytes"
	"github.com/magiconair/properties"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"sort"
)

func PropertiesToYaml(input []byte) ([]byte, error) {
	p, err := properties.Load(input, properties.UTF8)
	if nil != err {
		return nil, err
	}

	var propertyTree = Node{children: map[string]*Node{}}

	children := p.Map()
	keys := p.Keys()
	sort.Strings(keys)

	for _, key := range keys {
		propertyTree.deepSet(key, children[key])
	}

	toYaml, err := propertyTree.toYaml()
	if err == nil {
		return toYaml, nil
	} else {
		return nil, err
	}
}

func YamlToProperties(input []byte) ([]byte, error) {
	t := Node{}

	err := yaml.Unmarshal(input, &t)
	if err != nil {
		return nil, errors.Wrap(err, "yaml unmarshal")
	}
	var result bytes.Buffer
	err = t.toProperties(&result, "")
	return result.Bytes(), err
}
