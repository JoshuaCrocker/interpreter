package main

import (
	"errors"
	"regexp"
	"strconv"
)

const eof string = "EOF"
const number string = "NUMBER"
const operator string = "OPERATOR"

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
	getTokenParsers() []func(char string) token
	getNextToken() token
	eat(tokenType string) error
	Parse() string
}

func (i *interpreter) getTokenParsers() []func(char string) token {
	return []func(char string) token{
		// Number
		func(char string) token {
			rDigit, _ := regexp.Compile("[0-9]")

			if rDigit.MatchString(char) {
				return token{number, char}
			}

			return token{}
		},

		// Plus
		func(char string) token {
			if char == "+" {
				return token{operator, char}
			}

			return token{}
		},

		// Minus
		func(char string) token {
			if char == "-" {
				return token{operator, char}
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
	var t token
	blankToken := token{}

	for _, parser := range i.getTokenParsers() {
		t = parser(currentChar)

		if t != blankToken {
			i.pos++
			return t, nil
		}
	}

	return t, errors.New("Undefined token")
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
		return errors.New("Token Mismatch")
	}

	return nil
}

// Parse an arithmetic expression
func (i *interpreter) Parse() string {
	// Set the current token to the first token
	token, err := i.getNextToken()
	failOnError(err)

	i.currentToken = token

	left := i.currentToken
	err = i.eat(number)
	failOnError(err)

	op := i.currentToken
	err = i.eat(operator)
	failOnError(err)

	right := i.currentToken
	err = i.eat(number)
	failOnError(err)

	if op.Value == "+" {
		leftInt, _ := strconv.Atoi(left.Value.(string))
		rightInt, _ := strconv.Atoi(right.Value.(string))

		return strconv.Itoa(leftInt + rightInt)
	}

	if op.Value == "-" {
		leftInt, _ := strconv.Atoi(left.Value.(string))
		rightInt, _ := strconv.Atoi(right.Value.(string))

		return strconv.Itoa(leftInt - rightInt)
	}

	return ""
}
