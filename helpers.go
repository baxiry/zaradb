package dblite

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"runtime"
	"time"
)

func userDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

// shutdown
func Shutdown() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	// Shutdown grasefully
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	println("context")
	defer cancel()
	Server.Shutdown(ctx)

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
	cmd.Run()
	//Runs twice because sometimes pterodactyl servers needs a 2nd clear
}
