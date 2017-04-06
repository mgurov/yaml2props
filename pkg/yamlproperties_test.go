package pkg

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type exampleFile struct {
	Name     string
	Contents []byte
}

func readExample(fileName string) (result *exampleFile, err error) {
	result = &exampleFile{Name: fileName}
	result.Contents, err = ioutil.ReadFile(fileName)
	return
}

func checkExpectation(t *testing.T, from, to *exampleFile, conversion func([]byte) ([]byte, error)) {
	if actual, err := conversion(from.Contents); err != nil {
		t.Error(from.Name, "->", to.Name, err)
	} else {
		if string(actual) != string(to.Contents) {
			t.Errorf("%s->%s == \n->%s<-\n, want \n->%s<-\n", from.Name, to.Name, string(actual), string(to.Contents))
		}
	}

}

func Test2woWayConversions(t *testing.T) {
	yamlCases, err := filepath.Glob("../testcases/??-*.yaml")
	if nil != err {
		log.Fatal(err)
	}
	for _, yamlFileName := range yamlCases {
		yamlFile, err := readExample(yamlFileName)
		if err != nil {
			t.Fatal("Could not read yaml", err)
		}

		propsFile, err := readExample(yamlFileName[:len(yamlFileName)-4] + "properties")
		if err != nil {
			t.Fatal("Could not read props", err)
		}

		checkExpectation(t, yamlFile, propsFile, YamlToProperties)
		checkExpectation(t, propsFile, yamlFile, PropertiesToYaml)
	}
}

func TestInvalidYaml(t *testing.T) {
	sampleFile, err := readExample("../testcases/invalid-yaml.yaml")
	if err != nil {
		t.Fatal("Could not read file", err)
	}

	_, err = YamlToProperties(sampleFile.Contents)
	if nil == err {
		t.Error("Missed expected error")
	} else if !strings.Contains(err.Error(), "yaml unmarshal") {
		t.Error("Did not find expected error message within ", err)
	}
}

func TestMultilineYaml(t *testing.T) {
	sampleFile, err := readExample("../testcases/not-supported-multilines.yaml")
	if err != nil {
		t.Fatal("Could not read file", err)
	}

	result, err := YamlToProperties(sampleFile.Contents)
	if nil == err {
		t.Error("Missed expected error got instead:", result)
	} else if !strings.Contains(err.Error(), "multilined values not supported") {
		t.Error("Did not find expected error message within ", err)
	}
}

func TestCollectionYaml(t *testing.T) {
	sampleFile, err := readExample("../testcases/not-supported-collections.yaml")
	if err != nil {
		t.Fatal("Could not read file", err)
	}

	result, err := YamlToProperties(sampleFile.Contents)
	if nil == err {
		t.Error("Missed expected error got instead:", result)
	} else if !strings.Contains(err.Error(), "cannot unmarshal") {
		t.Error("Did not find expected error message within ", err)
	}
}

func TestMultilineProperties(t *testing.T) {
	sampleFile, err := readExample("../testcases/not-supported-multilines.properties")
	if err != nil {
		t.Fatal("Could not read file", err)
	}

	result, err := PropertiesToYaml(sampleFile.Contents)
	if nil == err {
		t.Error("Missed expected error got instead:", result)
	} else if !strings.Contains(err.Error(), "multilined values not supported") {
		t.Error("Did not find expected error message within ", err)
	}
}

func TestNotCompatibleProperties(t *testing.T) {
	sampleFile, err := readExample("../testcases/props-same-level.properties")
	if err != nil {
		t.Fatal("Could not read file", err)
	}

	_, err = PropertiesToYaml(sampleFile.Contents)
	if nil == err {
		t.Error("Missed expected error")
	} else if !strings.Contains(err.Error(), "invalid node .level1.level2") {
		t.Error("Did not find expected error message within ", err)
	}
}

func TestDeepSet(t *testing.T) {
	var actual Node
	actual.deepSet("a.b", "c")
	expected := Node{children: map[string]*Node{"a": {children: map[string]*Node{"b": {data: "c", children: map[string]*Node{}}}}}}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %s got %s", &expected, &actual)
	}
}

//TODO: docs mention no guarantee correct yaml
