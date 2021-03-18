package model

import (
	"strconv"

	"../stack"
)

type Calculator struct {
	dataStack        *stack.Stack
	operatorStack    *stack.Stack
	operatorPriority int //栈顶元素的优先级，0为最低，1为+-，2为*/
}

func NewCalculator() *Calculator {
	cal := Calculator{}
	cal.dataStack = stack.NewStack()
	cal.operatorStack = stack.NewStack()
	return &cal
}

func (cal *Calculator) calulate() {
	operator := cal.operatorStack.Pop()
	a := cal.dataStack.Pop()
	b := cal.dataStack.Pop()
	if operator == "*" {
		x, ok1 := a.(float64)
		y, ok2 := b.(float64)
		if ok1 && ok2 {
			cal.dataStack.Push(x * y)
		}
	}
	if operator == "/" {
		x, ok1 := a.(float64)
		y, ok2 := b.(float64)
		if ok1 && ok2 {
			cal.dataStack.Push(y / x)
		}
	}
	if operator == "+" {
		x, ok1 := a.(float64)
		y, ok2 := b.(float64)
		if ok1 && ok2 {
			cal.dataStack.Push(x + y)
		}
	}
	if operator == "-" {
		x, ok1 := a.(float64)
		y, ok2 := b.(float64)
		if ok1 && ok2 {
			cal.dataStack.Push(y - x)
		}
	}
}

func (cal *Calculator) CalculateFunc(result string, funcSlice []string, datas map[string]float64) {
	for _, p := range funcSlice {
		if p == "+" || p == "-" || p == "*" || p == "/" || p == "(" || p == ")" {
			if p == "+" || p == "-" {
				if cal.operatorPriority == 0 || cal.operatorPriority == 1 {
					cal.operatorStack.Push(p)
					cal.operatorPriority = 1
				}
				if cal.operatorPriority == 2 {
					cal.calulate()
					cal.operatorStack.Push(p)
					cal.operatorPriority = 1
				}
			}
			if p == "*" || p == "/" {
				cal.operatorStack.Push(p)
				cal.operatorPriority = 2
			}
			if p == "(" {
				cal.operatorStack.Push(p)
				cal.operatorPriority = 0
			}
			if p == ")" {
				for {
					if a, ok := cal.operatorStack.Top().(string); ok {
						if a == "(" {
							cal.operatorStack.Pop()
							if cal.operatorStack.IsEmpty() {
								cal.operatorPriority = 0
								break
							}
							a := cal.operatorStack.Top().(string)
							if a == "*" || a == "/" {
								cal.operatorPriority = 2
							}
							if a == "+" || a == "-" {
								cal.operatorPriority = 1
							}
							break
						}
					}
					cal.calulate()
				}
			}
		} else {
			if f, ok := strconv.ParseFloat(p, 64); ok == nil {
				cal.dataStack.Push(f)
			} else {
				cal.dataStack.Push(datas[p])
			}
		}
	}
	for !cal.operatorStack.IsEmpty() {
		cal.calulate()
	}
	a := cal.dataStack.Pop().(float64)
	datas[result] = a
}
