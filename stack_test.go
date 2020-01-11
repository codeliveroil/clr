package main

import (
	"testing"

	"github.com/codeliveroil/go-util/test"
)

func TestPush(t *testing.T) {
	s := &stack{}
	test.Compare(t, s.size, 0)

	s.push("a")
	test.Compare(t, s.toArray(), []colorCode{"a"})

	s.push("b")
	test.Compare(t, s.toArray(), []colorCode{"b", "a"})

	s.push("b")
	test.Compare(t, s.toArray(), []colorCode{"b", "b", "a"})
}

func TestPop(t *testing.T) {
	s := &stack{}

	s.push("a")
	s.push("b")
	s.push("b")

	verify := func(expE colorCode, expOK bool, expS []colorCode) {
		t.Helper()
		e, ok := s.pop()
		test.Compare(t, e, expE)
		test.Compare(t, ok, expOK)
		test.Compare(t, s.toArray(), expS)
	}

	verify("b", true, []colorCode{"b", "a"})
	verify("b", true, []colorCode{"a"})
	verify("a", true, []colorCode{})
	verify("", false, []colorCode{})
}

func TestPeek(t *testing.T) {
	s := &stack{}

	verify := func(expE colorCode, expOK bool, expS []colorCode) {
		t.Helper()
		e, ok := s.peek()
		test.Compare(t, e, expE)
		test.Compare(t, ok, expOK)
		test.Compare(t, s.toArray(), expS)
	}

	verify("", false, []colorCode{})

	s.push("a")
	verify("a", true, []colorCode{"a"})

	s.push("b")
	verify("b", true, []colorCode{"b", "a"})
}
