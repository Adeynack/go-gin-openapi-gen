package ginoascore

import (
	"errors"
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/getkin/kin-openapi/openapi3"
)

type Config struct {
	Specification *openapi3.Swagger
}

type Generation struct {
	Config *Config
	File   *jen.File
}

// Generate starts the code generation process.
func Generate(c *Config) (*Generation, error) {
	g := &Generation{
		c,
		jen.NewFile("api"),
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

func (g *Generation) generateSchemas(schemas map[string]*openapi3.SchemaRef) error {
	for schemaName, schema := range schemas {
		if err := g.generateSchema(schemaName, schema); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generation) generateSchema(schemaName string, schema *openapi3.SchemaRef) error {
	if schema.Value == nil {
		return fmt.Errorf("schema %q is unresolved schema (nil)", schemaName)
	}
	switch schema.Value.Type {
	case "object":
		g.generateObjectSchema(schemaName, schema.Value)
	default:
		statement := g.File.Type().Id(schemaName)
		g.addTypeToStatementFromSchemaRef(statement, schema)
	}

	return nil
}

func (g *Generation) generateObjectSchema(schemaName string, schema *openapi3.Schema) error {
	structProperties := make([]jen.Code, len(schema.Properties))
	i := 0
	for propName, prop := range schema.Properties {
		genProp := jen.Id(propName)
		err := g.addTypeToStatementFromSchemaRef(genProp, prop)
		if err != nil {
			return err
		}
		structProperties[i] = genProp
		i++
	}

	g.File.Type().Id(schemaName).Struct(structProperties...)
	return nil
}

func (g *Generation) addTypeToStatementFromSchemaRef(statement *jen.Statement, prop *openapi3.SchemaRef) error {
	if prop.Ref != "" {
		typeId, err := typeNameFromSchemaRef(prop.Ref)
		if err != nil {
			return err
		}
		statement.Qual(g.File.PackagePrefix, "*"+typeId)
		return nil
	}

	return g.addTypeToStatementFromSchema(statement, prop.Value)
}

func typeNameFromSchemaRef(ref string) (string, error) {
	var typeId string
	n, err := fmt.Sscanf(ref, "#/components/schemas/%s", &typeId)
	if err != nil || n != 1 {
		return "", fmt.Errorf(
			"could not extract schema name from ref %q (n: %d, err: %v)",
			ref, n, err)
	}
	return typeId, nil
}

func (g *Generation) addTypeToStatementFromSchema(s *jen.Statement, schema *openapi3.Schema) error {
	switch schema.Type {
	case "integer":
		s.Int() // todo: support `Format`!
	case "string":
		s.String() // todo: support `Format`?
	case "boolean":
		s.Bool()
	case "array":
		return g.completeArrayProperty(s, schema)
	default:
		return fmt.Errorf(
			"unable to generate property: unsupported schema type %q",
			schema.Type)
	}
	return nil
}

func (g *Generation) completeArrayProperty(statement *jen.Statement, schema *openapi3.Schema) error {
	if schema.Items.Ref == "" {
		return errors.New("arrays are only supported when their items is a ref")
	}
	typeName, err := typeNameFromSchemaRef(schema.Items.Ref)
	if err != nil {
		return err
	}
	statement.Index().Qual(g.File.PackagePrefix, `*`+typeName)
	return nil
}

/*
package api

type Book struct {
	id          *BookId		// todo: the JSON annotations with the OAS names
	name        string
	owner_id    *UserId		// todo: Transform the field names to CapitalizedCamelCase
	parent_list *BookList
}
type BookId int
type BookList struct {
	items *Book
}
type UserId int

 */
