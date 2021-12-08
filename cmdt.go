package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// root@158.247.195.235      e1J=xHytspxhhscx
func main() {
	deploy("158.247.195.235", "e1J=xHytspxhhscx", "test")
	os.Exit(1)

	saver("lilgo", "139.162.118.190", "d7ombot123")
}

// deploy deploys folder app to new remot server
func deploy(address, pass, foldname string) {

	//sshpass -p "password" scp -r /some/local/path user@example.com:/some/remote/path
	cmd := exec.Command("sshpass", "-p", pass, "scp", "-r", "./"+foldname, "root@"+address+":/root/") //.Output()

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(out))
	fmt.Println("Done")

}

// saver saves backup of deep client app folder
func saver(name, address, pass string) {
	for {
		//sshpass -p "password" scp -r user@example.com:/some/remote/path /some/local/path
		cmd := exec.Command("sshpass", "-p", pass, "scp", "-r", "root@"+address+":/root/"+name, "/root/") //.Output()

		out, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(out))
		fmt.Println("Done")

		time.Sleep(time.Second * 60 * 60)
	}
}
