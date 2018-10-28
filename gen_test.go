package ginoascore

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_LoadYamlSpecification(t *testing.T) {
	loader := openapi3.NewSwaggerLoader()
	assert.NotNil(t, loader)
	swagger, err := loader.LoadSwaggerFromYAMLData([]byte(openapiExampleFinances))
	if !assert.NoError(t, err) {
		return
	}
	conf := &Config{
		swagger,
	}
	assertSource(t, conf, func(t *testing.T, sa *SourceAssert) {
		sa.assertContainStruct("Book", []expectedStructProperty{
			{"id", "*BookId", ""},
			{"name", "string", ""},
			{"owner_id", "*UserId", ""},
		})
		sa.assertContainStruct("BookList", []expectedStructProperty{
			{"items", "[]*Book", ""},
		})
	})
}

func assertSource(t *testing.T, conf *Config, f func(t *testing.T, sourceAssert *SourceAssert)) {
	g, err := generate(conf)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, g) {
		return
	}
	sourceAssert := &SourceAssert{
		t:   t,
		src: g.file.GoString(),
	}
	f(t, sourceAssert)
}

type SourceAssert struct {
	t   *testing.T
	src string
}

type expectedStructProperty struct {
	name       string
	typ        string
	annotation string
}

const (
	rxSep    = `[\t ]+`
	rxSepOpt = `[\t ]*`
)

func (sa *SourceAssert) assertContainStruct(expectedStructName string, expectedProperties []expectedStructProperty) {
	if expectedStructName == "" {
		assert.Fail(sa.t, "an expected struct must have a name")
		return
	}
	rx := new(strings.Builder)
	rx.WriteString(`\ntype `)
	rx.WriteString(expectedStructName)
	rx.WriteString(` struct \{\n`)
	for i, p := range expectedProperties {
		rx.WriteString(rxSepOpt)
		if p.name == "" {
			assert.Failf(sa.t, "expecting struct %q: missing property name at index %d", expectedStructName, i)
			return
		}
		rx.WriteString(p.name)

		rx.WriteString(rxSep)
		typ := p.typ
		if typ == "" {
			assert.Failf(sa.t, "expecting struct %q: missing property type at index %d", expectedStructName, i)
			return
		}
		typ = strings.Replace(typ, "*", `\*`, -1)
		typ = strings.Replace(typ, "[", `\[`, -1)
		typ = strings.Replace(typ, "]", `\]`, -1)
		rx.WriteString(typ)

		if p.annotation != "" {
			rx.WriteString(rxSep)
			rx.WriteString("`")
			rx.WriteString(p.annotation)
			rx.WriteString("`")
		}
		rx.WriteString(`\n`)
	}
	rx.WriteString(`\}\n`)
	regEx := rx.String()
	assert.Regexp(sa.t, regEx, sa.src, "expected struct %q was not found or not as specified", expectedStructName)
}

const (
	openapiExampleFinances = `
openapi: 3.0.1

info:
  title: Finances
  version: 1.0.0

paths:

  /books:
    get:
      operationId: books.list
      description: |
        Get a list of all the books to which the current user has access.
        Users with 'admin' role will still the books to which they explicitlyæ
        have access to. Use the 'all=true' query parameter to see all of them.
      parameters:
      - $ref: '#/components/parameters/ListAll'
      responses:
        '200':
          description: List of books
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/BookList'

components:

  parameters:

    ListAll:
      name: list-all
      in: query
      description: |
        Force the complete list to be returned. This will make the operation fail
        if current user does not have necessary rights.
      required: false
      schema:
        type: boolean


  schemas:

    BookId:
      type: integer
      format: int64

    UserId:
      type: integer
      format: int64

    Book:
      type: object
      properties:
        id:
          $ref: '#/components/schemas/BookId'
        name:
          type: string
        owner_id:
          $ref: '#/components/schemas/UserId'

    BookList:
      type: object
      properties:
        items:
          type: array
          items:
            $ref: '#/components/schemas/Book'
`
)
