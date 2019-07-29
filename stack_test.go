package main

import (
	"testing"

	"github.com/codeliveroil/go-util/test"
)

func TestPush(t *testing.T) {
	s := &stack{}
	test.Compare(t, len(*s), 0)

	s.push("a")
	test.Compare(t, s, &stack{"a"})

	s.push("b")
	test.Compare(t, s, &stack{"b", "a"})

	s.push("b")
	test.Compare(t, s, &stack{"b", "b", "a"})
}

func TestPop(t *testing.T) {
	s := &stack{}

	s.push("a")
	s.push("b")
	s.push("b")

	verify := func(expE colorCode, expOK bool, expS *stack) {
		t.Helper()
		e, ok := s.pop()
		test.Compare(t, e, expE)
		test.Compare(t, ok, expOK)
		test.Compare(t, s, expS)
	}

	verify("b", true, &stack{"b", "a"})
	verify("b", true, &stack{"a"})
	verify("a", true, &stack{})
	verify("", false, &stack{})
}

func TestPeek(t *testing.T) {
	s := &stack{}

	verify := func(expE colorCode, expOK bool, expS *stack) {
		t.Helper()
		e, ok := s.peek()
		test.Compare(t, e, expE)
		test.Compare(t, ok, expOK)
		test.Compare(t, s, expS)
	}

	verify("", false, &stack{})

	s.push("a")
	verify("a", true, &stack{"a"})

	s.push("b")
	verify("b", true, &stack{"b", "a"})
}
