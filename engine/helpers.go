package engine

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
)

func int64ToBytes(n int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, n)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
		return []byte{}
	}
	return buf.Bytes()
}

func bytesToInt64(b []byte) int64 {
	var n int64
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &n)
	if err != nil {
		return 0
	}
	return n
}

func uint64ToBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b
}

func rootPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(usr.HomeDir, ".dbs") + "/" // slash
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

//end
