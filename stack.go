package main

import "fmt"

type colorCode string

type stack []colorCode

func (s *stack) pop() (colorCode, bool) {
	if len(*s) == 0 {
		return "", false
	}
	e := (*s)[0]
	*s = (*s)[1:]
	//printStack("Pop", e, *s)
	return e, true
}

func (s *stack) push(color colorCode) {
	*s = append([]colorCode{color}, *s...)
	//printStack("Push", color, *s)
}

func (s *stack) peek() (colorCode, bool) {
	if len(*s) == 0 {
		//printStack("Peek", "None", *s)
		return "", false
	}
	//printStack("Peek", (*s)[0], *s)

	return (*s)[0], true
}

func printStack(msg string, color colorCode, s stack) {
	fmt.Printf("%s: %sX\033[m -> ", msg, color)
	fmt.Printf("[")
	for _, e := range s {
		fmt.Printf("%sX\033[m,", e)
	}
	fmt.Printf("]\n")
}
