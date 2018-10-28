package main

import "flag"

var (
	inputSpecFile string
)

func main() {
	flag.StringVar(&inputSpecFile, "i", "", "path to the input OpenAPI Specification file (JSON or YAML)")

}
