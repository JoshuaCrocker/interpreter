package main

import (
	"errors"
	"regexp"
	"strconv"
)

const eof string = "EOF"
const number string = "NUMBER"
const plus string = "PLUS"

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
	GetNextToken() token
	Eat(tokenType string) error
	Parse() string
}

// GetNextToken from the input sentence
//
// Lexical analyser, breaks the input sentence down into tokens.
func (i *interpreter) GetNextToken() (token, error) {
	text := i.text

	if i.pos >= len(text) {
		return token{eof, ""}, nil
	}

	currentChar := string([]rune(text)[i.pos])

	rDigit, _ := regexp.Compile("[0-9]")

	if rDigit.MatchString(currentChar) {
		i.pos++
		return token{number, currentChar}, nil
	}

	if currentChar == "+" {
		i.pos++
		return token{plus, "+"}, nil
	}

	return token{}, errors.New("Undefined token")
}

func (i *interpreter) Eat(tokenType string) error {
	if i.currentToken.Type == tokenType {
		next, err := i.GetNextToken()
		if nil != err {
			panic(err)
		}

		i.currentToken = next
	} else {
		return errors.New("Token Mismatch")
	}

	return nil
}

func (i *interpreter) Parse() string {
	// Set the current token to the first token
	token, err := i.GetNextToken()
	failOnError(err)

	i.currentToken = token

	left := i.currentToken
	err = i.Eat(number)
	failOnError(err)

	op := i.currentToken
	err = i.Eat(plus)
	failOnError(err)

	right := i.currentToken
	err = i.Eat(number)
	failOnError(err)

	if op.Value == "+" {
		leftInt, _ := strconv.Atoi(left.Value.(string))
		rightInt, _ := strconv.Atoi(right.Value.(string))

		return strconv.Itoa(leftInt + rightInt)
	}

	return ""
}
