package main

import "fmt"

func failOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	i := interpreter{"9+8", 0, token{}}
	fmt.Printf("9+8= %s", i.Parse())
}
