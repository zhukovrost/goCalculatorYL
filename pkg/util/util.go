package util

import (
	"fmt"
	"strings"
	"unicode"
)

var precedence = map[rune]int{
	'+': 1,
	'-': 1,
	'*': 2,
	'/': 2,
}

func isOperator(c rune) bool {
	_, exists := precedence[c]
	return exists
}

func higherPrecedence(op1, op2 rune) bool {
	return precedence[op1] >= precedence[op2]
}

func ToPostfix(expression string) ([]string, error) {
	var stack []rune
	var output []string
	var numberBuffer strings.Builder

	for _, token := range expression {
		switch {
		case unicode.IsDigit(token):
			numberBuffer.WriteRune(token)
		case token == '(':
			if numberBuffer.Len() > 0 {
				output = append(output, numberBuffer.String())
				numberBuffer.Reset()
			}
			stack = append(stack, token)
		case token == ')':
			if numberBuffer.Len() > 0 {
				output = append(output, numberBuffer.String())
				numberBuffer.Reset()
			}
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				output = append(output, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, fmt.Errorf("mismatched parentheses")
			}
			stack = stack[:len(stack)-1]
		case isOperator(token):
			if numberBuffer.Len() > 0 {
				output = append(output, numberBuffer.String())
				numberBuffer.Reset()
			}
			for len(stack) > 0 && isOperator(stack[len(stack)-1]) && higherPrecedence(stack[len(stack)-1], token) {
				output = append(output, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		}
	}

	if numberBuffer.Len() > 0 {
		output = append(output, numberBuffer.String())
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == '(' {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		output = append(output, string(stack[len(stack)-1]))
		stack = stack[:len(stack)-1]
	}

	return output, nil
}
