package pkg

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"sort"
	"strings"
)

type Node struct {
	data     string
	children map[string]*Node
}

// see gopkg.in/yaml.v2
func (e *Node) UnmarshalYAML(unmarshal func(interface{}) error) error {

	if err := unmarshal(&e.data); err == nil {
		//stringable -> scalar
		return nil
	}

	return unmarshal(&e.children)
}

func (e *Node) toYaml() (b []byte, err error) {
	buf := new(bytes.Buffer)
	err = e.toYamlInternal(buf, "", -1, "")
	return buf.Bytes(), err
}

func (e *Node) toYamlInternal(to *bytes.Buffer, name string, level int, path string) error {
	if e.data != "" && len(e.children) > 0 {
		return errors.New("Error converting to yaml: invalid node " + path + " cannot contain scalar value and children same node")
	}

	if strings.Contains(e.data, "\n") {
		return errors.New("Error converting to yaml: invalid node " + path + " multilined values not supported")
	}

	if level >= 0 {
		whitespace := strings.Repeat("  ", level)
		fmt.Fprintf(to, "%s%s:", whitespace, name)
		if "" != e.data {
			fmt.Fprint(to, " ", e.data)
		}
		to.WriteString("\n")
	}
	for _, k := range e.sortedKeys() {
		if err := e.children[k].toYamlInternal(to, k, level+1, path+"."+k); err != nil {
			return err
		}
	}
	return nil
}

func (e *Node) String() string {
	return fmt.Sprintf("%s(%s)", e.data, e.children)
}

func (e *Node) toProperties(w *bytes.Buffer, path string) error {
	if e.data != "" {
		if strings.Contains(e.data, "\n") {
			return errors.Errorf("Error converting to yaml property %s: multilined values not supported", path)
		}
		fmt.Fprintf(w, "%s=%s\n", path, e.data)
	}

	for _, key := range e.sortedKeys() {
		childPath := key
		if "" != path {
			childPath = path + "." + childPath
		}
		err := e.children[key].toProperties(w, childPath)
		if nil != err {
			return err
		}
	}
	return nil
}

func (e *Node) sortedKeys() []string {
	keys := make([]string, len(e.children))
	i := 0
	for key, _ := range e.children {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

func (d *Node) deepSet(path string, value string) {
	node := d
	pathParts := strings.Split(path, ".")
	for i, p := range pathParts {
		if nil == node.children {
			node.children = map[string]*Node{}
		}
		child, ok := node.children[p]
		if !ok {
			child = &Node{children: map[string]*Node{}}
			node.children[p] = child
		}
		node = child
		if i == len(pathParts)-1 {
			node.data = value
		}
	}
}
