package main

import "fmt"

type colorCode string

type node struct {
	next  *node
	color colorCode
}

type stack struct {
	top  *node
	size int
}

func (s *stack) push(color colorCode) {
	n := &node{
		color: color,
	}

	s.size++

	if s.top == nil {
		s.top = n
		return
	}

	n.next = s.top
	s.top = n
}

func (s *stack) pop() (colorCode, bool) {
	if s.top == nil {
		return "", false
	}

	c := s.top.color
	s.size--

	if s.top.next == nil {
		s.top = nil
		return c, true
	}

	s.top = s.top.next

	return c, true
}

func (s *stack) peek() (colorCode, bool) {
	if s.top == nil {
		return "", false
	}

	return s.top.color, true
}

func printStack(msg string, color colorCode, s stack) {
	fmt.Printf("%s: %sX\033[m -> ", msg, color)
	fmt.Printf("[")

	n := s.top
	for n != nil {
		fmt.Printf("%sX\033[m,", n.color)
		n = n.next
	}

	fmt.Printf("]\n")
}

func (s *stack) toArray() (arr []colorCode) {
	n := s.top
	for n != nil {
		arr = append(arr, n.color)
		n = n.next
	}
	if len(arr) == 0 {
		// for reflect.DeepEqual in test cases. Otherwise []colorCode{}
		// won't equal what is returned when the stack is empty
		return []colorCode{}
	}
	return arr
}
