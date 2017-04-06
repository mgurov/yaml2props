package main

import (
	"flag"
	"fmt"
	"github.com/mgurov/yaml2props/pkg"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	os.Exit(cli(os.Args, os.Stdin))
}

func cli(args []string, stdin io.Reader) int {

	appName := args[0]
	flagSet := flag.NewFlagSet(appName, flag.ContinueOnError)
	reverse := flagSet.Bool("r", false, "reverse")
	showVersion := flagSet.Bool("version", false, "show version and exit")

	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: \n%s [<flags>] [<input file name, stdin if omitted>]\nflags:\n", appName)
		flagSet.PrintDefaults()

	}
	err := flagSet.Parse(args[1:])
	if nil != err {
		log.Println("parsing args", err)
		return 1
	}

	if *showVersion {
		fmt.Println("0.1")
		return 0
	}

	var inputStream io.Reader

	if flagSet.NArg() == 1 {
		fileName := flagSet.Arg(0)
		f, err := os.Open(fileName)
		if err != nil {
			log.Println("Could not read", fileName, err)
			return 1
		}
		defer f.Close()
		inputStream = f
	} else if flagSet.NArg() > 1 {
		log.Println("expect one or zero arguments")
		flagSet.Usage()
		return 1
	} else {
		inputStream = stdin
	}

	input, err := ioutil.ReadAll(inputStream)
	if err != nil {
		log.Println("Could not read stdin", err)
		return 1
	}

	var output []byte
	if *reverse {
		output, err = pkg.PropertiesToYaml(input)
	} else {
		output, err = pkg.YamlToProperties(input)
	}

	if err != nil {
		log.Println(err)
		return 1
	}
	os.Stdout.Write(output)
	return 0
}
