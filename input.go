package main

import (
	"fmt"
)

func main() {
	var in string
	fmt.Scanf("%s\n", &in)
	fmt.Println(in)
}

func hi() {
	fmt.Errorf("noting%s", "hei")
}
