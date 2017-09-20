package generatejsschema

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	jsschema "github.com/lestrrat/go-jsschema"
)

func testEqual(t *testing.T, got, expected interface{}) {
	t.Helper()
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Error(diff)
	}
}

func TestFromArgs(t *testing.T) {
	g := Generator{}
	err := g.FromArgs([]string{"testdata/types.go"})
	if err != nil {
		t.Fatal(err)
	}

	schema := g.Schema

	t.Run("description", func(t *testing.T) {
		testEqual(t, schema.Definitions["User"].Description, "User is a data type for user")
	})

	t.Run("enum", func(t *testing.T) {
		testEqual(t, schema.Definitions["UserStatus"].Enum, []interface{}{"active", "inactive"})
	})

	t.Run("array", func(t *testing.T) {
		testEqual(t, schema.Definitions["User"].Properties["Items"].Type, jsschema.PrimitiveTypes{jsschema.ArrayType})
		testEqual(t, schema.Definitions["User"].Properties["Items"].Items.Schemas[0].Reference, "#/definitions/Item")
	})
}
