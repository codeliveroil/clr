package main

import (
	"regexp"
	"testing"

	"github.com/codeliveroil/go-util/test"
)

func TestParse(t *testing.T) {
	cases := []struct {
		input string
		exp   *filter
	}{
		{"kobe~red", &filter{
			exp:           regexp.MustCompile("kobe"),
			code:          "\033[38;5;1m",
			highlightLine: false,
		}},
		{"shaq~12,blue", &filter{
			exp:           regexp.MustCompile("shaq"),
			code:          "\033[38;5;12m\033[48;5;4m",
			highlightLine: false,
		}},
		{"jordan [[:digit:]]~default,yellow~line", &filter{
			exp:           regexp.MustCompile("jordan [[:digit:]]"),
			code:          "\033[39m\033[48;5;3m",
			highlightLine: true,
		}},
	}

	args := make([]string, len(cases))
	for i, c := range cases {
		args[i] = c.input
	}

	filters := parse(args)

	test.Compare(t, filters != nil, true)
	test.Compare(t, len(filters), len(cases))

	for i, c := range cases {
		got := filters[i]
		want := c.exp
		test.Compare(t, got, want)
	}
}

func TestColor(t *testing.T) {
	var got string

	//Test boundaries and interleaving of colors
	got = color("this is a test string", parse([]string{
		"this is a test string~red",
		"is a test~yellow,white",
		"a~blue",
		"g~253",
		"n~default,pink~line",
	}))

	test.Compare(t, got, "\033[38;5;1m\033[39m\033[48;5;5mthis \033[38;5;3m\033[48;5;15mis \033[38;5;4ma\033[m\033[38;5;3m\033[48;5;15m test\033[m\033[39m\033[48;5;5m strin\033[38;5;253mg\033[m")

	//Test regex
	got = color("Hello, W0rld 2019!", parse([]string{
		"Hello~red",
		"[[:digit:]]~blue,pink",
	}))

	test.Compare(t, got, "\033[38;5;1mHello\033[m\033[m, W\033[38;5;4m\033[48;5;5m0\033[m\033[mrld \033[38;5;4m\033[48;5;5m2\033[m\033[m\033[38;5;4m\033[48;5;5m0\033[m\033[m\033[38;5;4m\033[48;5;5m1\033[m\033[m\033[38;5;4m\033[48;5;5m9\033[m\033[m!\033[m")
}
