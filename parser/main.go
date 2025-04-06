package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/stu-k/go/parser/parse"
)

type Data interface {
	Type() string
	Value() any
	String() string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		fmt.Printf("parsing: \"%v\"\n", input)
		result, err := parse.Parse(input)
		if err != nil {
			fmt.Printf("error parsing input: %v\n", err)
			continue
		}

		fmt.Printf("result: %+v\n", result)
	}
}
