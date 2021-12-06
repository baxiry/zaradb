package save

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {

	data := readFile("main.go")
	writeFile("main.go2", data)
}

func readFile(path string) string {
	data := ""
	f, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		data += scanner.Text() + "\n"

		fmt.Println(scanner.Text())
		// do somethin with this data
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return data
}
func writeFile(path, data string) {
	f, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(data)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("done")
}

//////////////////////////////////////////////
