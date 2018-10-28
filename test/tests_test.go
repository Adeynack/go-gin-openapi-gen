package test

import (
	"github.com/adeynack/go-gin-openapi-gen"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func Test_Finances(t *testing.T) {
	testFromYaml(t, "finances")
}

func testFromYaml(t *testing.T, specificationName string) {
	loader := openapi3.NewSwaggerLoader()
	assert.NotNil(t, loader)
	specificationFileContent, err := ioutil.ReadFile(specificationName + ".yaml")
	if !assert.NoError(t, err) {
		return
	}
	swagger, err := loader.LoadSwaggerFromYAMLData(specificationFileContent)
	if !assert.NoError(t, err) {
		return
	}
	conf := &ginoascore.Config{
		Specification: swagger,
	}
	g, err := ginoascore.Generate(conf)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, g) {
		return
	}
	generatedSource := g.File.GoString()

	generatedSourceFile := specificationName + "_generated.txt"
	err = ioutil.WriteFile(generatedSourceFile, []byte(generatedSource), os.ModePerm)
	if !assert.NoError(t, err) {
		return
	}
	t.Logf(
		"Generated code was outputed to %q (this file is GIT-ignored and will not be commited)",
		generatedSourceFile)

	expectedSource, err := ioutil.ReadFile(specificationName + "_expected.txt")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, string(expectedSource), generatedSource)
}
