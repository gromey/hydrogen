# Hydrogen

![https://img.shields.io/github/v/tag/gromey/hydrogen](https://img.shields.io/github/v/tag/gromey/hydrogen)
![https://img.shields.io/github/license/gromey/hydrogen](https://img.shields.io/github/license/gromey/hydrogen)

## Overview

`hydrogen` is a library designed to get names, values and pointers to structure fields.

## Installation

`hydrogen` can be installed like any other Go library through `go get`:

```sh
  go get github.com/gromey/hydrogen@latest
```

## Getting Started

```go
package main

import (
	"fmt"

	"github.com/gromey/hydrogen"
)

type Type struct {
	A string // will be returned with the 'A' name
	B string `tag:"b"` // will be returned with the 'b' name
	C string `tag:"-"` // will be skipped
	D string `tag:"d"`
}

func main() {
	tag := hydrogen.New("tag")

	s := new(Type)
	s.A = "fa"
	s.B = "fb"
	s.C = "fc"
	s.D = ""

	fls, err := tag.Fields(s, false)
	if err != nil {
		panic(err)
	}

	names := fls.Names()
	fmt.Printf("names: %v\n", names)
	// names: [A b d]

	values := fls.Values()
	fmt.Printf("values: %v\n", values)
	// values: [fa fb ]

	pointers := fls.Pointers()
	fmt.Println("pointers:", &s.A == pointers[0], &s.B == pointers[1], &s.D == pointers[2])
	// pointers: true, true, true

	// s.D will be omitted because it's empty
	fls, err = tag.Fields(s, true)
	if err != nil {
		panic(err)
	}

	names = fls.Names()
	fmt.Printf("names: %v\n", names)
	// names: [A b]

	values = fls.Values()
	fmt.Printf("values: %v\n", values)
	// values: [fa fb]

	pointers = fls.Pointers()
	fmt.Println("pointers:", &s.A == pointers[0], &s.B == pointers[1])
	// pointers: true, true
}

```
