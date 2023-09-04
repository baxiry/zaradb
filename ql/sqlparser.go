package parser

import (
	"fmt"
	"strings"
)

func parse(query string) {

	tokens := make([]string, 0)
	args := make([]string, 0)
	words := make([]string, 0)

	words = strings.Split(query, " ")

	for _, w := range words {
		if w == "" {
			continue
		}
		if w[:1] == "_" {
			tokens = append(tokens, w)
		} else {
			args = append(args, w)
		}
	}

	println()
	fmt.Println("words: ", words, len(words))
	println()
	fmt.Println("tokens: ", tokens, len(tokens))
	println()
	fmt.Println("args: ", args, len(args))
}
