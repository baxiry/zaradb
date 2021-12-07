package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {

	for {
		saver("lilgo", "139.162.118.190", "d7ombot123")
		time.Sleep(time.Second * 60 * 60)
	}

}
func saver(name, address, pass string) {

	//sshpass -p "password" scp -r user@example.com:/some/remote/path /some/local/path
	cmd := exec.Command("sshpass", "-p", pass, "scp", "-r", "root@"+address+":/root/"+name, "/root/") //.Output()

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(out))
	fmt.Println("Done")

}
