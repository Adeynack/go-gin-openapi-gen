package gen

import (
	"github.com/dave/jennifer/jen"
	"github.com/getkin/kin-openapi/openapi3"
)

// Config of a code generation.
type Config struct {
	Specification   *openapi3.Swagger
	OutputDirectory string
}

// Generation state.
type Generation struct {
	Config     *Config
	File       *jen.File
	SchemaInfo map[string]*SchemaInfo
}

// SchemaInfo represents the metainformation of a generated schema.
type SchemaInfo struct {
	Schema   *openapi3.SchemaRef
	TypeRef  bool
	IsStruct bool
}

// Generate starts the code generation process.
func Generate(c *Config) (*Generation, error) {
	g := &Generation{
		c,
		jen.NewFile("api"),
		make(map[string]*SchemaInfo),
	}

	err := g.generateComponents(&g.Config.Specification.Components)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Generation) generateComponents(components *openapi3.Components) error {
	err := g.generateSchemas(components.Schemas)
	if err != nil {
		return err
	}
	return nil
}
