package hydrogen

import (
	"reflect"
	"testing"
)

func equal(t *testing.T, exp, got interface{}) {
	if !reflect.DeepEqual(exp, got) {
		t.Fatalf("Not equal:\nexp: %v\ngot: %v", exp, got)
	}
}

func Test_isEmptyValue(t *testing.T) {
	a := 77
	var b any

	var tests = []struct {
		name   string
		value  any
		expect bool
	}{
		{
			name:   "no empty bool",
			value:  true,
			expect: false,
		},
		{
			name:   "empty bool",
			value:  false,
			expect: true,
		},
		{
			name:   "no empty int",
			value:  1,
			expect: false,
		},
		{
			name:   "empty int",
			value:  0,
			expect: true,
		},
		{
			name:   "np empty float",
			value:  1.1,
			expect: false,
		},
		{
			name:   "empty float",
			value:  0.0,
			expect: true,
		},
		{
			name:   "no empty string",
			value:  "a",
			expect: false,
		},
		{
			name:   "empty string",
			value:  "",
			expect: true,
		},
		{
			name:   "no empty pointer",
			value:  &struct{}{},
			expect: false,
		},
		{
			name:   "empty pointer",
			value:  (*struct{})(nil),
			expect: true,
		},
		{
			name:   "no empty slice",
			value:  []int{1},
			expect: false,
		},
		{
			name:   "empty slice",
			value:  []int{},
			expect: true,
		},
		{
			name:   "no empty interface",
			value:  a,
			expect: false,
		},
		{
			name:   "empty interface",
			value:  b,
			expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			equal(t, tt.expect, isEmptyValue(reflect.ValueOf(tt.value)))
		})
	}
}
