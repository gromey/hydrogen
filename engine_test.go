package hydrogen

import (
	"errors"
	"testing"
)

var tag = New("tag")

type simple struct {
	A string
	B int
	C float64
	D int
	e string
}

type nested struct {
	A string
	B int
	C float64 `tag:"name"`
	simple
	S *simple
}

type skip struct {
	A string `tag:"-"`
	B int    `tag:"-"`
	C float64
	*simple
	S simple
}

type pointer struct {
	A *int
}

type deepEmbed struct {
	Foo nested `tag:"foo"`
	Bar string `tag:"bar"`
}

var (
	ns = 12345

	s1 = simple{A: "test", B: 0, C: 3.14, D: 28, e: "not exported"}
	s2 = nested{A: "foo", B: 0, C: 3.14, simple: s1, S: &s1}
	s3 = skip{A: "foo", B: 0, C: 3.14, simple: &s1, S: s1}
	s4 = pointer{}
	s5 = deepEmbed{Foo: s2, Bar: "text"}
)

func TestEngine_Fields(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		omitempty bool
		fields    []string
		expect    Fields
		err       error
	}{
		{
			name:  "not a pointer",
			input: s1,
			err:   errors.New("the input value is not a pointer"),
		},
		{
			name:  "not a struct",
			input: &ns,
			err:   errors.New("the input value is not a struct"),
		},
		{
			name:  "all fields",
			input: &s1,
			expect: Fields{
				{Name: "A", Value: "test", Addr: &s1.A},
				{Name: "B", Value: 0, Addr: &s1.B},
				{Name: "C", Value: 3.14, Addr: &s1.C},
				{Name: "D", Value: 28, Addr: &s1.D},
			},
		},
		{
			name:      "all fields with omitempty",
			input:     &s1,
			omitempty: true,
			expect: Fields{
				{Name: "A", Value: "test", Addr: &s1.A},
				{Name: "C", Value: 3.14, Addr: &s1.C},
				{Name: "D", Value: 28, Addr: &s1.D},
			},
		},
		{
			name:   "selected fields",
			input:  &s1,
			fields: []string{"A", "C"},
			expect: Fields{
				{Name: "A", Value: "test", Addr: &s1.A},
				{Name: "C", Value: 3.14, Addr: &s1.C},
			},
		},
		{
			name:      "selected fields with omitempty",
			input:     &s1,
			omitempty: true,
			fields:    []string{"A", "B"},
			expect: Fields{
				{Name: "A", Value: "test", Addr: &s1.A},
			},
		},
		{
			name:  "all fields with nested",
			input: &s2,
			expect: Fields{
				{Name: "A", Value: "foo", Addr: &s2.A},
				{Name: "B", Value: 0, Addr: &s2.B},
				{Name: "name", Value: 3.14, Addr: &s2.C},
				{Name: "simple_A", Value: "test", Addr: &s2.simple.A},
				{Name: "simple_B", Value: 0, Addr: &s2.simple.B},
				{Name: "simple_C", Value: 3.14, Addr: &s2.simple.C},
				{Name: "simple_D", Value: 28, Addr: &s2.D},
				{Name: "S_A", Value: "test", Addr: &s2.S.A},
				{Name: "S_B", Value: 0, Addr: &s2.S.B},
				{Name: "S_C", Value: 3.14, Addr: &s2.S.C},
				{Name: "S_D", Value: 28, Addr: &s2.S.D},
			},
		},
		{
			name:      "all fields with nested with omitempty",
			input:     &s2,
			omitempty: true,
			expect: Fields{
				{Name: "A", Value: "foo", Addr: &s2.A},
				{Name: "name", Value: 3.14, Addr: &s2.C},
				{Name: "simple_A", Value: "test", Addr: &s2.simple.A},
				{Name: "simple_C", Value: 3.14, Addr: &s2.simple.C},
				{Name: "simple_D", Value: 28, Addr: &s2.D},
				{Name: "S_A", Value: "test", Addr: &s2.S.A},
				{Name: "S_C", Value: 3.14, Addr: &s2.S.C},
				{Name: "S_D", Value: 28, Addr: &s2.S.D},
			},
		},
		{
			name:   "selected fields with nested",
			input:  &s2,
			fields: []string{"A", "B", "simple_A", "simple_B", "S_A", "S_B"},
			expect: Fields{
				{Name: "A", Value: "foo", Addr: &s2.A},
				{Name: "B", Value: 0, Addr: &s2.B},
				{Name: "simple_A", Value: "test", Addr: &s2.simple.A},
				{Name: "simple_B", Value: 0, Addr: &s2.simple.B},
				{Name: "S_A", Value: "test", Addr: &s2.S.A},
				{Name: "S_B", Value: 0, Addr: &s2.S.B},
			},
		},
		{
			name:      "selected fields with nested with omitempty",
			input:     &s2,
			omitempty: true,
			fields:    []string{"A", "B", "simple_A", "simple_B", "S_A", "S_B"},
			expect: Fields{
				{Name: "A", Value: "foo", Addr: &s2.A},
				{Name: "simple_A", Value: "test", Addr: &s2.simple.A},
				{Name: "S_A", Value: "test", Addr: &s2.S.A},
			},
		},
		{
			name:  "all fields with skip",
			input: &s3,
			expect: Fields{
				{Name: "C", Value: 3.14, Addr: &s3.C},
				{Name: "simple_A", Value: "test", Addr: &s3.simple.A},
				{Name: "simple_B", Value: 0, Addr: &s3.simple.B},
				{Name: "simple_C", Value: 3.14, Addr: &s3.simple.C},
				{Name: "simple_D", Value: 28, Addr: &s3.D},
				{Name: "S_A", Value: "test", Addr: &s3.S.A},
				{Name: "S_B", Value: 0, Addr: &s3.S.B},
				{Name: "S_C", Value: 3.14, Addr: &s3.S.C},
				{Name: "S_D", Value: 28, Addr: &s3.S.D},
			},
		},
		{
			name:      "all fields with skip with omitempty",
			input:     &s3,
			omitempty: true,
			expect: Fields{
				{Name: "C", Value: 3.14, Addr: &s3.C},
				{Name: "simple_A", Value: "test", Addr: &s3.simple.A},
				{Name: "simple_C", Value: 3.14, Addr: &s3.simple.C},
				{Name: "simple_D", Value: 28, Addr: &s3.D},
				{Name: "S_A", Value: "test", Addr: &s3.S.A},
				{Name: "S_C", Value: 3.14, Addr: &s3.S.C},
				{Name: "S_D", Value: 28, Addr: &s3.S.D},
			},
		},
		{
			name:   "selected fields with skip",
			input:  &s3,
			fields: []string{"A", "C", "simple_A", "simple_C", "S_A", "S_C"},
			expect: Fields{
				{Name: "C", Value: 3.14, Addr: &s3.C},
				{Name: "simple_A", Value: "test", Addr: &s3.simple.A},
				{Name: "simple_C", Value: 3.14, Addr: &s3.simple.C},
				{Name: "S_A", Value: "test", Addr: &s3.S.A},
				{Name: "S_C", Value: 3.14, Addr: &s3.S.C},
			},
		},
		{
			name:      "selected fields with skip with omitempty",
			input:     &s3,
			omitempty: true,
			fields:    []string{"A", "B", "simple_A", "simple_B", "S_A", "S_B"},
			expect: Fields{
				{Name: "simple_A", Value: "test", Addr: &s3.simple.A},
				{Name: "S_A", Value: "test", Addr: &s3.S.A},
			},
		},
		{
			name:  "nil pointer field",
			input: &s4,
			expect: Fields{
				{Name: "A", Value: (*int)(nil), Addr: &s4.A},
			},
		},
		{
			name:      "nil pointer field",
			input:     &s4,
			omitempty: true,
			expect:    Fields{},
		},
		{
			name:  "deep embedding",
			input: &s5,
			expect: Fields{
				{Name: "foo_A", Value: "foo", Addr: &s5.Foo.A},
				{Name: "foo_B", Value: 0, Addr: &s5.Foo.B},
				{Name: "foo_name", Value: 3.14, Addr: &s5.Foo.C},
				{Name: "foo_simple_A", Value: "test", Addr: &s5.Foo.simple.A},
				{Name: "foo_simple_B", Value: 0, Addr: &s5.Foo.simple.B},
				{Name: "foo_simple_C", Value: 3.14, Addr: &s5.Foo.simple.C},
				{Name: "foo_simple_D", Value: 28, Addr: &s5.Foo.D},
				{Name: "foo_S_A", Value: "test", Addr: &s5.Foo.S.A},
				{Name: "foo_S_B", Value: 0, Addr: &s5.Foo.S.B},
				{Name: "foo_S_C", Value: 3.14, Addr: &s5.Foo.S.C},
				{Name: "foo_S_D", Value: 28, Addr: &s5.Foo.S.D},
				{Name: "bar", Value: "text", Addr: &s5.Bar},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := tag.Fields(tt.input, tt.omitempty, tt.fields...)
			if tt.err != nil {
				equal(t, tt.err.Error(), err.Error())
				return
			} else {
				equal(t, nil, err)
			}
			equal(t, len(tt.expect), len(fields))
			for i, f := range fields {
				equal(t, tt.expect[i].Name, f.Name)
				equal(t, tt.expect[i].Value, f.Value)
				equal(t, tt.expect[i].Addr, f.Addr)
			}
		})
	}
}

func TestFields_Names(t *testing.T) {
	tests := []struct {
		name   string
		input  any
		expect []string
	}{
		{
			name:   "get names",
			input:  &s2,
			expect: []string{"A", "B", "name", "simple_A", "simple_B", "simple_C", "simple_D", "S_A", "S_B", "S_C", "S_D"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := tag.Fields(tt.input, false)
			equal(t, nil, err)

			names := fields.Names()
			equal(t, tt.expect, names)
		})
	}
}

func TestFields_Values(t *testing.T) {
	tests := []struct {
		name   string
		input  any
		expect []any
	}{
		{
			name:   "get values",
			input:  &s2,
			expect: []any{"foo", 0, 3.14, "test", 0, 3.14, 28, "test", 0, 3.14, 28},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := tag.Fields(tt.input, false)
			equal(t, nil, err)

			values := fields.Values()
			equal(t, tt.expect, values)
		})
	}
}

func TestFields_Pointers(t *testing.T) {
	tests := []struct {
		name   string
		input  any
		expect []any
	}{
		{
			name:   "get pointers",
			input:  &s2,
			expect: []any{&s2.A, &s2.B, &s2.C, &s2.simple.A, &s2.simple.B, &s2.simple.C, &s2.D, &s2.S.A, &s2.S.B, &s2.S.C, &s2.S.D},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := tag.Fields(tt.input, false)
			equal(t, nil, err)

			pointers := fields.Pointers()
			equal(t, tt.expect, pointers)
		})
	}
}
