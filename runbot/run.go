package main

import (
	"bufio"
	//"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {

	tokenFile := "token.txt"
	var bot string
	if len(os.Args) > 1 {
		bot = os.Args[1]
	}
	fmt.Println(bot)

	file, err := os.Open(tokenFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	args := make(map[string]string)
	line := make([]string, 2)

	for scanner.Scan() {

		line = strings.Split(scanner.Text(), " ")
		args[line[0]] = strings.Trim(line[1], " ")

	}

	cmd := exec.Command(os.Args[1], args[os.Args[1]])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}

//bot := flag.String("b", "bot-0", "use -b <name_Bot> to change default bot")
//tok := flag.String("t", "token.txt", "use -t <name_token_file> to change default tokensFile")
//flag.Parse()
