package ginoascore

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
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
	g, err := generate(conf)
	if !assert.NoError(t, err) {
		return
	}

	assert.NotNil(t, g)

	r := g.file.GoString()

	assert.Contains(t, r, "type Book struct {")
	assert.Contains(t, r, "type BookList struct {")

	println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	println(r)
	println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
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
        parent_list:
          $ref: '#/components/schemas/BookList'

    BookList:
      type: object
      properties:
        items:
          type: array
          items:
            $ref: '#/components/schemas/Book'
`
)
