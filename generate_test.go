package generatejsschema

import (
	"fmt"
	"testing"
)

func TestFromArgs(t *testing.T) {
	g := Generator{}
	err := g.FromArgs([]string{"testdata/types.go"})
	if err != nil {
		t.Fatal(err)
	}

	if got, expected := g.Schema.Definitions["User"].Description, "User is a data type for user"; got != expected {
		t.Errorf("User/Description: got %v != %v", got, expected)
	}

	if got, expected := fmt.Sprint(g.Schema.Definitions["UserStatus"].Enum), "[active inactive]"; got != expected {
		t.Errorf("UserStatus/Enum: got %v != %v", got, expected)
	}
}
