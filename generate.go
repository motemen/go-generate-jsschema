package generatejsschema

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/loader"
	"log"
	"reflect"
	"strings"

	jsschema "github.com/lestrrat/go-jsschema"
	"github.com/pkg/errors"
)

type Generator struct {
	Schema  *jsschema.Schema
	program *loader.Program
}

func (g *Generator) FromArgs(args []string) error {
	conf := loader.Config{
		ParserMode: parser.ParseComments,
	}
	_, err := conf.FromArgs(args, false)
	if err != nil {
		return err
	}

	g.program, err = conf.Load()
	if err != nil {
		return err
	}

	g.Schema = jsschema.New()
	g.Schema.SchemaRef = jsschema.SchemaURL
	g.Schema.Definitions = map[string]*jsschema.Schema{}
	g.Schema.AdditionalItems = &jsschema.AdditionalItems{}
	g.Schema.AdditionalProperties = &jsschema.AdditionalProperties{}

	for _, pkg := range g.program.InitialPackages() {
		pkgScope := pkg.Pkg.Scope()
		for _, name := range pkgScope.Names() {
			obj := pkgScope.Lookup(name)
			if !obj.Exported() {
				continue
			}

			tn, ok := obj.(*types.TypeName)
			if !ok {
				continue
			}

			var err error
			g.Schema.Definitions[tn.Name()], err = g.processType(tn.Type().Underlying(), obj)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Generator) processType(typ types.Type, obj types.Object) (*jsschema.Schema, error) {
	debugf("processType: %s", typ)

	switch typ := typ.(type) {
	case *types.Array:
		debugf("not implemented")
	case *types.Basic:
		return g.processBasicType(typ)
	case *types.Chan:
		debugf("not implemented")
	case *types.Interface:
		debugf("not implemented")
	case *types.Map:
		debugf("not implemented")
	case *types.Named:
		debugf("not implemented")
	case *types.Pointer:
		debugf("not implemented")
	case *types.Signature:
		debugf("not implemented")
	case *types.Slice:
		debugf("not implemented")
	case *types.Struct:
		return g.processStructType(typ, obj)
	case *types.Tuple:
		debugf("not implemented")
	}

	return nil, errors.Errorf("not implemented for: %s", obj)
}

func (g *Generator) processStructType(st *types.Struct, obj types.Object) (*jsschema.Schema, error) {
	schema := jsschema.New()
	schema.Properties = map[string]*jsschema.Schema{}
	schema.Description = g.docStringAtPos(obj.Pos())

	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		if !field.Exported() {
			continue
		}

		propSchema, err := g.processType(field.Type(), field)
		if err != nil {
			return nil, err
		}

		tag := reflect.StructTag(st.Tag(i)).Get("json")
		if tag == "-" {
			continue
		}

		name, opts := parseTag(tag)
		if name == "" {
			name = field.Name()
		}

		propSchema.Description = g.docStringAtPos(field.Pos())

		schema.Properties[name] = propSchema

		if !opts.Contains("omitempty") {
			schema.Required = append(schema.Required, name)
		}
	}

	return schema, nil
}

func (g *Generator) processBasicType(bt *types.Basic) (*jsschema.Schema, error) {
	schema := jsschema.New()
	switch bt.Kind() {
	case types.Bool:
		schema.Type = jsschema.PrimitiveTypes{jsschema.BooleanType}
		return schema, nil

	// case types.Byte: == types.Uint8
	case types.Complex128:
	case types.Complex64:

	case types.Float32:
	case types.Float64:

	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
		schema.Type = jsschema.PrimitiveTypes{jsschema.IntegerType}
		return schema, nil

	case types.Invalid:

	// case types.Rune: == types.Int32

	case types.String:
		schema.Type = jsschema.PrimitiveTypes{jsschema.StringType}
		return schema, nil

	case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
		// TODO: minValue?
		schema.Type = jsschema.PrimitiveTypes{jsschema.IntegerType}
		return schema, nil

	case types.Uintptr:
	case types.UnsafePointer:
	case types.UntypedBool:
	case types.UntypedComplex:
	case types.UntypedFloat:
	case types.UntypedInt:
	case types.UntypedNil:
	case types.UntypedRune:
	case types.UntypedString:
	}

	return nil, errors.Errorf("not implemented for type: %s", bt)
}

func (g *Generator) docStringAtPos(pos token.Pos) string {
	var comment *ast.CommentGroup

	_, nodes, _ := g.program.PathEnclosingInterval(pos, pos)
	for _, node := range nodes {
		switch node := node.(type) {
		case *ast.Field:
			comment = node.Doc
			goto last
		case *ast.GenDecl:
			comment = node.Doc
			goto last
		}
	}

last:
	if comment != nil {
		return strings.TrimSpace(comment.Text())
	}

	return ""
}

func debugf(format string, args ...interface{}) {
	log.Printf("debug: "+format, args...)
}
