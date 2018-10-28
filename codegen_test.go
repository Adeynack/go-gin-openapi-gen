package ginoascore

import (
	. "github.com/dave/jennifer/jen"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	f := NewFile("foo")

	typeBook := f.Type().Id("Book").Struct(
		Id("Id").Int64(),
		Id("Name").String(),
		Id("OwnerId").Int64(),
	)
	assert.NotNil(t, typeBook)

	s := f.GoString()
	println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	println(s)
	println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")

	assert.NotEmpty(t, s, f)
}
