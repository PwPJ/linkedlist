package main

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestLinkedListProperties(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 320 //successful test cases run per property
	parameters.MaxSize = 80             //maxsize for generated data

	properties := gopter.NewProperties(parameters)

	// test the insert and get property
	properties.Property("test insert and get", prop.ForAll(
		func(a int, b int) bool {
			l := New()
			l.Insert(0, a)
			l.Insert(1, b)
			va, _ := l.Get(0)
			vb, _ := l.Get(1)
			return va == a && vb == b
		},
		gen.Int(),
		gen.Int(),
	))

	properties.Property("length after insert and Remove", prop.ForAll(
		func(a int) bool {
			l := New()
			l.Insert(0, a)
			success := l.Remove(0)
			return success && l.length == 0
		},
		gen.Int(),
	))

	properties.Property("keeping the order the same after adding many items", prop.ForAll(
		func(a []int) bool {
			l := New()
			for i, v := range a {
				l.Insert(uint(i), v)
			}
			for i, v := range a {
				value, _ := l.Get(uint(i))
				if value != v {
					return false
				}
			}
			return true
		},
		gen.SliceOf(gen.Int()),
	))

	properties.TestingRun(t)
}
