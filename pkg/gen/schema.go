package gen

import (
	"errors"
	"fmt"
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/getkin/kin-openapi/openapi3"
)

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
	schemaInfo := new(SchemaInfo)
	g.SchemaInfo[schemaName] = schemaInfo
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
	}

	statement.Struct(structProperties...)
	return nil
}

func (g *Generation) addTypeToStatementFromSchemaRef(s *jen.Statement, prop *openapi3.SchemaRef) (*jen.Statement, error) {
	if prop.Ref != "" {
		typeID, err := typeNameFromSchemaRef(prop.Ref)
		if err != nil {
			return nil, err
		}
		// todo: For struct types, generate pointer qualified identifier. (*Book instead of Book)
		s = s.Qual(g.File.PackagePrefix, typeID)
		return s, nil
	}

	return g.addTypeToStatementFromSchema(s, prop.Value)
}

func typeNameFromSchemaRef(ref string) (string, error) {
	var typeID string
	n, err := fmt.Sscanf(ref, "#/components/schemas/%s", &typeID)
	if err != nil || n != 1 {
		return "", fmt.Errorf(
			"could not extract schema name from ref %q (n: %d, err: %v)",
			ref, n, err)
	}
	goTypeID := toGoFieldName(typeID)
	return goTypeID, nil
}

func (g *Generation) addTypeToStatementFromSchema(s *jen.Statement, schema *openapi3.Schema) (*jen.Statement, error) {
	// https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.2.md#data-types
	switch schema.Type {
	case "integer":
		return g.addTypeIntToStatementFromSchema(s, schema)
	case "number":
		return g.addTypeNumberToStatementFromSchema(s, schema)
	case "string":
		return g.addTypeStringToStatementFromSchema(s, schema)
	case "boolean":
		return s.Bool(), nil
	case "array":
		return g.completeArrayProperty(s, schema)
	default:
		return nil, fmt.Errorf(
			"unable to generate property: unsupported schema type %q",
			schema.Type)
	}
}

func (g *Generation) addTypeIntToStatementFromSchema(s *jen.Statement, schema *openapi3.Schema) (*jen.Statement, error) {
	// https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.2.md#data-types
	switch schema.Format {
	case "", "int32":
		return s.Int(), nil
	case "int64":
		return s.Int64(), nil
	default:
		return nil, fmt.Errorf("integer format %q is not recognized", schema.Format)
	}
}

func (g *Generation) addTypeNumberToStatementFromSchema(s *jen.Statement, schema *openapi3.Schema) (*jen.Statement, error) {
	// https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.2.md#data-types
	switch schema.Format {
	case "", "float":
		return s.Float32(), nil
	case "double":
		return s.Float64(), nil
	default:
		return nil, fmt.Errorf("number format %q is not recognized", schema.Format)
	}
}

func (g *Generation) addTypeStringToStatementFromSchema(s *jen.Statement, schema *openapi3.Schema) (*jen.Statement, error) {
	// https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.2.md#data-types
	switch schema.Format {
	case "", "password":
		return s.String(), nil
	case "byte":
		return s.Index().Byte(), nil
	case "binary":
		return nil, errors.New(`string format "binary" is not yet supported`)
	case "date", "date-time":
		return s.Qual("time", "Time"), nil
	default:
		return nil, fmt.Errorf("string format %q is not recognized", schema.Format)
	}
}

func (g *Generation) completeArrayProperty(s *jen.Statement, schema *openapi3.Schema) (*jen.Statement, error) {
	if schema.Items.Ref == "" {
		return nil, errors.New("arrays are only supported when their items is a ref")
	}
	typeName, err := typeNameFromSchemaRef(schema.Items.Ref)
	if err != nil {
		return nil, err
	}
	s = s.Index()
	// todo: For struct types, generate pointer qualified identifier. (*Book instead of Book)
	s = s.Qual(g.File.PackagePrefix, typeName)
	return s, nil
}
