package test

import (
	"github.com/adeynack/go-gin-openapi-gen/pkg/gen"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func Test_Finances(t *testing.T) {
	testFromYaml(t, "finances", nil)
}

func testFromYaml(
	t *testing.T,
	specificationName string,
	swaggerTestFn func(t *testing.T, swagger *openapi3.Swagger),
) (generation *gen.Generation, swagger *openapi3.Swagger) {
	loader := openapi3.NewSwaggerLoader()
	assert.NotNil(t, loader)
	specificationFileContent, err := ioutil.ReadFile(specificationName + ".yaml")
	if !assert.NoError(t, err) {
		return
	}
	swagger, err = loader.LoadSwaggerFromYAMLData(specificationFileContent)
	if !assert.NoError(t, err) {
		return
	}
	if swaggerTestFn != nil {
		swaggerTestFn(t, swagger)
	}
	conf := &gen.Config{
		Specification: swagger,
	}
	generation, err = gen.Generate(conf)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, generation) {
		return
	}
	generatedSource := generation.File.GoString()

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
	return
}
