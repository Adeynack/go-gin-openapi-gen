package gen

import (
	"errors"
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/getkin/kin-openapi/openapi3"
	"sort"
)

type Config struct {
	Specification   *openapi3.Swagger
	OutputDirectory string
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
	// Sort schemas by name
	schemaNames := make([]string, 0, len(schemas))
	for n := range schemas {
		schemaNames = append(schemaNames, n)
	}
	sort.Strings(schemaNames)
	// Create schemas in alphabetical order
	for _, schemaName := range schemaNames {
		if err := g.generateSchema(schemaName, schemas[schemaName]); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generation) generateSchema(schemaName string, schema *openapi3.SchemaRef) error {
	if schema.Value == nil {
		return fmt.Errorf("schema %q is unresolved schema (nil)", schemaName)
	}
	g.File.Commentf("// %s implements OpenAPI element at #/components/schemas/%s", schemaName, schemaName)
	g.File.Commentf(schema.Value.Description)
	statement := g.File.Type().Id(schemaName)
	switch schema.Value.Type {
	case "object":
		g.addStructToStatementFromSchema(statement, schema.Value)
	default:
		g.addTypeToStatementFromSchemaRef(statement, schema)
	}

	return nil
}

func (g *Generation) addStructToStatementFromSchema(statement *jen.Statement, schema *openapi3.Schema) error {
	// Sort properties by name
	propertyNames := make([]string, 0, len(schema.Properties))
	for n := range schema.Properties {
		propertyNames = append(propertyNames, n)
	}
	sort.Strings(propertyNames)
	// Create properties in alphabetical order
	structProperties := make([]jen.Code, len(schema.Properties))
	for i, propName := range propertyNames {
		goPropertyName := toGoFieldName(propName)
		genProp := jen.Id(goPropertyName)
		genProp, err := g.addTypeToStatementFromSchemaRef(genProp, schema.Properties[propName])
		if err != nil {
			return err
		}
		genProp.Tag(map[string]string{"json": propName})
		structProperties[i] = genProp
		i++
	}

	statement.Struct(structProperties...)
	return nil
}

func (g *Generation) addTypeToStatementFromSchemaRef(s *jen.Statement, prop *openapi3.SchemaRef) (*jen.Statement, error) {
	if prop.Ref != "" {
		typeId, err := typeNameFromSchemaRef(prop.Ref)
		if err != nil {
			return nil, err
		}
		s = s.Qual(g.File.PackagePrefix, "*"+typeId)
		return s, nil
	}

	return g.addTypeToStatementFromSchema(s, prop.Value)
}

func typeNameFromSchemaRef(ref string) (string, error) {
	var typeId string
	n, err := fmt.Sscanf(ref, "#/components/schemas/%s", &typeId)
	if err != nil || n != 1 {
		return "", fmt.Errorf(
			"could not extract schema name from ref %q (n: %d, err: %v)",
			ref, n, err)
	}
	goTypeId := toGoFieldName(typeId)
	return goTypeId, nil
}

func (g *Generation) addTypeToStatementFromSchema(s *jen.Statement, schema *openapi3.Schema) (*jen.Statement, error) {
	switch schema.Type {
	case "integer":
		s = s.Int() // todo: support `Format`!
	case "string":
		s = s.String() // todo: support `Format`?
	case "boolean":
		s = s.Bool()
	case "array":
		return g.completeArrayProperty(s, schema)
	default:
		return nil, fmt.Errorf(
			"unable to generate property: unsupported schema type %q",
			schema.Type)
	}
	return s, nil
}

func (g *Generation) completeArrayProperty(s *jen.Statement, schema *openapi3.Schema) (*jen.Statement, error) {
	if schema.Items.Ref == "" {
		return nil, errors.New("arrays are only supported when their items is a ref")
	}
	typeName, err := typeNameFromSchemaRef(schema.Items.Ref)
	if err != nil {
		return nil, err
	}
	s = s.Index().Qual(g.File.PackagePrefix, `*`+typeName)
	return s, nil
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
