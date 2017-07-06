package generatejsschema

import (
	"encoding/json"
	"testing"
)

func TestFromArgs(t *testing.T) {
	g := Generator{}
	err := g.FromArgs([]string{"testdata/types.go"})
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.Marshal(g.Schema)
	t.Log(string(b))
}
