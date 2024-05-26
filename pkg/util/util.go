package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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

	for i, token := range expression {
		switch {
		case unicode.IsDigit(token) || token == '.': // Support for decimal numbers
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
		case unicode.IsSpace(token):
			if numberBuffer.Len() > 0 && i > 0 && unicode.IsDigit(rune(expression[i-1])) {
				output = append(output, numberBuffer.String())
				numberBuffer.Reset()
			}
		default:
			return nil, fmt.Errorf("invalid character: %c", token)
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

func GenerateId() string {
	return strconv.FormatInt(time.Now().UnixNano()%1000000, 10)
}
