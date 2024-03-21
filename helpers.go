package store

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"runtime"
)

// returns file path separator depand on os
func Slash() string {
	if runtime.GOOS == "windows" {
		return "\\"
	}
	return "/"
}

func rootPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(usr.HomeDir, ".dbs") + "/"
}

func ClearScreen() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
	//cmd.Run()
	//Runs twice because sometimes pterodactyl servers needs a 2nd clear
}

func SysNotify() {

	var c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	s := <-c
	fmt.Println("Got signal:", s)

}

// PathExist check if path exists
func PathExist(subPath string) bool {
	_, err := os.Stat(rootPath() + subPath)
	return os.IsNotExist(err)
}

// ListDir show all directories in path
func ListDir(path string) {
	dbs, err := os.ReadDir(rootPath() + path)
	if err != nil {
		fmt.Println(err)
	}

	dirs := 0
	for _, dir := range dbs {
		if dir.IsDir() && string(dir.Name()[0]) != "." {
			dirs++
			print(dir.Name(), " ")
		}
	}
	if dirs > 0 {
		println()
		return
	}
	println(path, "is impty")
}
