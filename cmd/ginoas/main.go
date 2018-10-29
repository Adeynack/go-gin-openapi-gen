package main

import (
	"flag"
	"fmt"
	"github.com/adeynack/go-gin-openapi-gen/pkg/gen"
	"github.com/getkin/kin-openapi/openapi3"
	"io/ioutil"
	"os"
)

var (
	flagInputSpecFile   string
	flagOutputDirectory string
)

func main() {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		fmt.Printf("%v", err)
		os.Exit(1)
	}()

	flag.StringVar(&flagInputSpecFile, "i", "", "path to the input OpenAPI Specification file (JSON or YAML)")
	flag.StringVar(&flagOutputDirectory, "o", "", "directory in which the files are to be generated")
	flag.Parse()

	if flagInputSpecFile == "" {
		fmt.Fprintln(os.Stderr, "No input file provided (-i PATH)")
		os.Exit(1)
	}

	swagger := loadSwaggerFromYaml(flagInputSpecFile)

	conf := &gen.Config{
		Specification:   swagger,
		OutputDirectory: flagOutputDirectory,
	}

	generation, err := gen.Generate(conf)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "Generation completed: %#v", generation)
	err = generation.File.Render(os.Stdout)
	if err != nil {
		panic(fmt.Errorf("error rendering generated source: %v", err))
	}
}

func loadSwaggerFromYaml(specFile string) *openapi3.Swagger {
	bytes, err := ioutil.ReadFile(specFile)
	if err != nil {
		panic(fmt.Errorf("error reading file %q: %v", specFile, err))
	}
	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromYAMLData(bytes)
	if err != nil {
		panic(fmt.Errorf("error loading OpenAPI specification from file %q: %v", specFile, err))
	}
	return swagger
}
