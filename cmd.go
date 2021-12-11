package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	saver("lilgo", "139.162.118.190", "d7ombot123")
	os.Exit(1)
	deploy("test", "158.247.195.235", "e1J=xHytspxhhscx")

	for {
		time.Sleep(time.Second * 60 * 60)
	}

}

// deploy deploys folder app to new remot server
func deploy(foldname, address, pass string) {

	//sshpass -p "password" scp -r /some/local/path user@example.com:/some/remote/path
	cmd := exec.Command("sshpass", "-p", pass, "scp", "-r", "./"+foldname, "root@"+address+":/root/") //.Output()

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(out))
	fmt.Println("Done")

}

func saver(filename, address, pass string) {

	//sshpass -p "password" scp -r user@example.com:/some/remote/path /some/local/path
	cmd := exec.Command("sshpass", "-p", pass, "scp", "-r", "root@"+address+":/root/"+filename, "/root/") //.Output()

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(out))
	fmt.Println("Done")

}
