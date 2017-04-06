package main

import (
	"log"
	"strings"
)

func ExampleFromYamlStdin() {
	yaml := "a: 1\nb: 2"
	if err := cli([]string{"app"}, strings.NewReader(yaml)); err != 0 {
		log.Println("No success execution code")
	}

	// Output:
	// a=1
	// b=2
}

func ExampleFromPropertiesStdin() {
	yaml := "a: 1\nb: 2"
	if err := cli([]string{"app", "-r"}, strings.NewReader(yaml)); err != 0 {
		log.Println("No success execution code")
	}

	// Output:
	// a: 1
	// b: 2
}

func ExampleFromFile() {
	if err := cli([]string{"app", "../../testcases/01-simple.yaml"}, strings.NewReader("")); err != 0 {
		log.Println("No success execution code")
	}

	// Output:
	// name=MyName
	// subvalue.how-much=1.1
}
