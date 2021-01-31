package main

import (
	"errors"
	"regexp"
	"strconv"
)

const eof string = "EOF"
const number string = "NUMBER"
const plus string = "PLUS"
const minus string = "MINUS"
const multiply string = "MULTIPLY"
const divide string = "DIVIDE"

type token struct {
	Type  string
	Value interface{}
}

type interpreter struct {
	text         string
	pos          int
	currentToken token
}

type iinterpreter interface {
	getTokenParsers() []func(char string, text string, pos int) token
	getNextToken() token
	eat(tokenType string) error
	Parse() string
}

func (i *interpreter) getTokenParsers() []func(char string, text string, pos int) token {
	return []func(char string, text string, pos int) token{
		// Number
		func(char string, text string, pos int) token {
			rDigit, _ := regexp.Compile("[0-9]")
			output := ""

			for rDigit.MatchString(char) {
				output += char
				pos++

				if pos < len(text) {
					char = string([]rune(text)[pos])
				} else {
					char = ""
				}
			}

			if output != "" {
				return token{number, output}
			}

			return token{}
		},

		// Plus
		func(char string, text string, pos int) token {
			if char == "+" {
				return token{plus, char}
			}

			return token{}
		},

		// Minus
		func(char string, text string, pos int) token {
			if char == "-" {
				return token{minus, char}
			}

			return token{}
		},

		// Multiply
		func(char string, text string, pos int) token {
			if char == "*" {
				return token{multiply, char}
			}

			return token{}
		},

		// Divide
		func(char string, text string, pos int) token {
			if char == "/" {
				return token{divide, char}
			}

			return token{}
		},
	}
}

// GetNextToken from the input sentence
//
// Lexical analyser, breaks the input sentence down into tokens by getting the
// next character within the sentence, and iterates through the known parsers
// to get the next token.
//
// Fails if no parser functions matches the given input.
func (i *interpreter) getNextToken() (token, error) {
	text := i.text

	if i.pos >= len(text) {
		return token{eof, ""}, nil
	}

	currentChar := string([]rune(text)[i.pos])
	for currentChar == " " {
		i.pos++
		currentChar = string([]rune(text)[i.pos])
	}
	var t token
	blankToken := token{}

	for _, parser := range i.getTokenParsers() {
		t = parser(currentChar, i.text, i.pos)

		if t != blankToken {
			n := 0
			for n < len(t.Value.(string)) {
				i.pos++
				n++
			}

			return t, nil
		}
	}

	return t, errors.New("Undefined token: " + currentChar)
}

// Eat the next token within the input. This method fails if the next token
// is not of the expected type.
func (i *interpreter) eat(tokenType string) error {
	if i.currentToken.Type == tokenType {
		next, err := i.getNextToken()
		if nil != err {
			panic(err)
		}

		i.currentToken = next
	} else {
		return errors.New("Token Mismatch (got " + i.currentToken.Type + " expected " + tokenType + ")")
	}

	return nil
}

// Parse an arithmetic expression
func (i *interpreter) Parse() string {
	var left, op, right token
	var output int = 0

	// Set the current token to the first token
	token, err := i.getNextToken()
	failOnError(err)

	i.currentToken = token

	left = i.currentToken
	err = i.eat(number)
	failOnError(err)

	output, _ = strconv.Atoi(left.Value.(string))

	for i.currentToken.Type != eof {
		op = i.currentToken
		var err error = nil
		if op.Type == plus {
			err = i.eat(plus)
		} else if op.Type == minus {
			err = i.eat(minus)
		} else if op.Type == multiply {
			err = i.eat(multiply)
		} else if op.Type == divide {
			err = i.eat(divide)
		}
		failOnError(err)

		right = i.currentToken
		err = i.eat(number)
		failOnError(err)

		if op.Value == "+" {
			rightInt, _ := strconv.Atoi(right.Value.(string))
			output += rightInt
		} else if op.Value == "-" {
			rightInt, _ := strconv.Atoi(right.Value.(string))
			output -= rightInt
		} else if op.Value == "*" {
			rightInt, _ := strconv.Atoi(right.Value.(string))
			output *= rightInt
		} else if op.Value == "/" {
			rightInt, _ := strconv.Atoi(right.Value.(string))
			output /= rightInt
		}
	}

	return strconv.Itoa(output)
}
