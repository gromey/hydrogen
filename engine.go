package hydrogen

import (
	"errors"
	"reflect"
	"strings"
	"sync"
)

type Engine interface {
	Fields(v any, omitempty bool, fieldNames ...string) (Fields, error)
}

func New(name string) Engine {
	return &engine{name: name}
}

type engine struct {
	name string
}

func (e *engine) Fields(v any, omitempty bool, fieldNames ...string) (Fields, error) {
	s := e.newState()
	defer statePool.Put(s)

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return nil, errors.New("the input value is not a pointer")
	}

	var err error
	if rv, err = valueFromPtr(rv); err != nil {
		return nil, err
	}

	fs := s.cachedFields(rv.Type())
	return fs.capture(rv, omitempty, "", fieldNames...), nil
}

type state struct {
	*engine
}

var statePool sync.Pool

func (e *engine) newState() *state {
	if p := statePool.Get(); p != nil {
		return p.(*state)
	}
	return &state{engine: e}
}

// field represents a single field found in a struct.
type field struct {
	index    int
	name     string
	embedded structFields
}

type structFields []field

var fieldCache sync.Map // map[reflect.Type]structFields

// cachedFields is like typeFields but uses a cache to avoid repeated work.
func (e *engine) cachedFields(t reflect.Type) structFields {
	if c, ok := fieldCache.Load(t); ok {
		return c.(structFields)
	}
	c, _ := fieldCache.LoadOrStore(t, e.typeFields(t))
	return c.(structFields)
}

// typeFields returns a list of fields for the given type.
func (e *engine) typeFields(t reflect.Type) structFields {
	fs := make(structFields, 0, t.NumField())

	// Scan t for fields.
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		ft := sf.Type

		f := field{
			index: i,
			name:  sf.Name,
		}

		if ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
		}

		if sf.Anonymous {
			// Ignore embedded fields of unexported non-struct types.
			if !sf.IsExported() && ft.Kind() != reflect.Struct {
				continue
			}

			// Do not ignore embedded fields of unexported struct types since they may have exported fields.
			f.embedded = e.cachedFields(ft)
			fs = append(fs, f)

			continue
		} else if !sf.IsExported() {
			// Ignore unexported non-embedded fields.
			continue
		}

		if tag, ok := sf.Tag.Lookup(e.name); ok {
			// Ignore the field if the tag has a skip value.
			if tag == "-" {
				continue
			}

			if n := strings.IndexByte(tag, ','); n < 0 && len(tag) > 0 {
				f.name = strings.ReplaceAll(tag, " ", "_")
			} else if n > 0 {
				f.name = strings.ReplaceAll(tag[:n], " ", "_")
			}
		}

		if ft.Kind() == reflect.Struct {
			f.embedded = e.cachedFields(ft)
			fs = append(fs, f)

			continue
		}

		fs = append(fs, f)
	}

	return fs
}

func (f *structFields) capture(v reflect.Value, omitempty bool, parent string, fieldNames ...string) Fields {
	fs := make(Fields, 0, len(*f))

	for _, fld := range *f {
		rv := v.Field(fld.index)

		if rv.Kind() == reflect.Pointer {
			rv = rv.Elem()
		}

		if isEmptyValue(rv) && omitempty {
			continue
		}

		name := parent + fld.name

		if fld.embedded != nil {
			fs = append(fs, fld.embedded.capture(rv, omitempty, name+"_", fieldNames...)...)
			continue
		}

		if len(fieldNames) != 0 && ignore(name, fieldNames...) {
			continue
		}

		fs = append(fs, Field{
			Name:  name,
			Value: v.Field(fld.index).Interface(),
			Addr:  v.Field(fld.index).Addr().Interface(),
		})
	}

	return fs
}

type Field struct {
	Name  string
	Value any
	Addr  any
}

type Fields []Field

// Names returns a slice of field names.
func (f Fields) Names() []string {
	names := make([]string, 0, len(f))
	for _, fld := range f {
		names = append(names, fld.Name)
	}
	return names
}

// Values returns a slice of field values.
func (f Fields) Values() []any {
	values := make([]any, 0, len(f))
	for _, fld := range f {
		values = append(values, fld.Value)
	}
	return values
}

// Pointers returns a slice of field pointers.
func (f Fields) Pointers() []any {
	pointers := make([]any, 0, len(f))
	for _, fld := range f {
		pointers = append(pointers, fld.Addr)
	}
	return pointers
}
