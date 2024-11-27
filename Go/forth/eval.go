//go:build !solution

package main

import (
	"errors"
	"strconv"
	"strings"
)

type Stack struct {
	stack []int
}

func (s *Stack) Push(v int) {
	s.stack = append(s.stack, v)
}

func (s *Stack) Pop() (int, error) {
	if len(s.stack) == 0 {
		return 0, errors.New("stack is empty")
	}
	v := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return v, nil
}

func (s *Stack) Top() (int, error) {
	if len(s.stack) == 0 {
		return 0, errors.New("stack is empty")
	}
	return s.stack[len(s.stack)-1], nil
}

func (s *Stack) Top2() (int, error) {
	if len(s.stack) < 2 {
		return 0, errors.New("not enough elements")
	}
	return s.stack[len(s.stack)-2], nil
}

func (s *Stack) Len() int {
	return len(s.stack)
}

type Evaluator struct {
	m     map[string]func(stack *Stack) error
	stack Stack
}

func NewEvaluator() *Evaluator {
	var ev = Evaluator{make(map[string]func(stack *Stack) error), Stack{make([]int, 0)}}
	ev.m["+"] = Add
	ev.m["-"] = Sub
	ev.m["*"] = Mult
	ev.m["/"] = Div
	ev.m["dup"] = Dup
	ev.m["over"] = Over
	ev.m["drop"] = Drop
	ev.m["swap"] = Swap
	return &ev
}

func Add(stack *Stack) error {
	var int1, _ = stack.Pop()
	var int2, err = stack.Pop()
	if err == nil {
		stack.Push(int1 + int2)
	} else {
		return err
	}
	return err
}

func Sub(stack *Stack) error {
	var int1, _ = stack.Pop()
	var int2, err = stack.Pop()
	if err == nil {
		stack.Push(int2 - int1)
	} else {
		return err
	}
	return err
}

func Mult(stack *Stack) error {
	var int1, _ = stack.Pop()
	var int2, err = stack.Pop()
	if err == nil {
		stack.Push(int1 * int2)
	} else {
		return err
	}
	return err
}

func Div(stack *Stack) error {
	var int1, _ = stack.Pop()
	var int2, err = stack.Pop()
	if err == nil {
		if int1 == 0 {
			return errors.New("division by zero")
		}
		stack.Push(int2 / int1)
	} else {
		return err
	}
	return err
}

func Dup(stack *Stack) error {
	var int1, err = stack.Top()
	if err == nil {
		stack.Push(int1)
	} else {
		return err
	}
	return err
}

func Drop(stack *Stack) error {
	var _, err = stack.Pop()
	return err
}

func Swap(stack *Stack) error {
	var int1, _ = stack.Pop()
	var int2, err = stack.Pop()
	if err == nil {
		stack.Push(int1)
		stack.Push(int2)
	}
	return err
}

func Over(stack *Stack) error {
	var int1, err = stack.Top2()
	if err == nil {
		stack.Push(int1)
	}
	return err
}

func (e *Evaluator) Process(row string) ([]int, error) {
	row = strings.ToLower(row)
	splitted := strings.Split(row, " ")
	var err error
	if len(splitted) >= 4 && splitted[0] == ":" && splitted[len(splitted)-1] == ";" {
		splitted = splitted[1 : len(splitted)-1]
		name := splitted[0]
		var _, err2 = strconv.Atoi(name)
		if err2 == nil {
			return nil, errors.New("can not redefine number")
		}
		funcNames := splitted[1:]
		functions := make([]func(stack *Stack) error, 0)
		for _, funcName := range funcNames {
			function, ok := e.m[funcName]
			if ok {
				functions = append(functions, function)
			} else {
				var number int
				number, err2 = strconv.Atoi(funcName)
				if err2 == nil {
					functions = append(functions, func(stack *Stack) error {
						stack.Push(number)
						return nil
					})
				}
			}
		}
		e.m[name] = func(stack *Stack) error {
			for _, function := range functions {
				err = function(stack)
			}
			return err
		}
	} else {
		i := 0
		for i < len(splitted) {
			arg := splitted[i]
			var number, err2 = strconv.Atoi(arg)
			if err2 == nil {
				e.stack.Push(number)
			} else {
				f, ok := e.m[arg]
				if ok {
					err = f(&e.stack)
					if err != nil {
						return nil, err
					}
				} else {
					return nil, errors.New("no such word")
				}
			}
			i++
		}
	}
	return e.stack.stack, err
}
