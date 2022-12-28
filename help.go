package helper

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/melbahja/goph"
)

// to scure app read pass form seprite file
func Getpass() string {
	data, err := ioutil.ReadFile(".mypass")
	if err != nil {
		return err.Error()
	}
	psw := string(data)
	return psw[:len(psw)-1]
}

// file copies a single file from src to dst
func Copylocalfile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

// Copydir copies local botline directory
// this is copies a whole directory recursively
func CopyDir(src string, dst string) error {
	dst = dst + "-bot"
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		fmt.Println("err: os.stat")
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		fmt.Println("err: os.makeall")
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		fmt.Println("err: ioutil.readdir")
		return err
	}
	for _, fd := range fds {
		srcfp := filepath.Join(src, fd.Name())
		dstfp := filepath.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {

				fmt.Println("err: recoursive 1")
				fmt.Println(err)
			}
		} else {
			if err = Copylocalfile(srcfp, dstfp); err != nil {
				fmt.Println("err: recoursive 2")
				fmt.Println(err)
			}
		}
	}

	// creat a new file that containe client info,
	clientinfo, err := os.Create(dst + "/" + dst + ".info")
	if err != nil {
		fmt.Println("creat file info when copping dir")
		return err
	}
	defer clientinfo.Close()
	clientinfo.WriteString(dst)
	return nil
}

// randsleep sleep program 100 to 1000 millisecond
func Randsleep() {
	rand.Seed(time.Now().UnixMicro())
	t := rand.Intn(900)
	time.Sleep(time.Millisecond * time.Duration(t+100))
}

// todo test zip function
//
//	zipfile.zip and clientname
func RemoteZip(sshclient *goph.Client, outfile, dir string) error {
	// zip the client bot app
	cmd, err := sshclient.Command("zip", "-r", outfile+".zip", dir)
	if err != nil {
		return err
	}
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// removeitem remove item string from list and return new list
func RemoveItem(item string, list []string) []string {
	newlist := make([]string, 0)
	for _, v := range list {
		if item != v {
			newlist = append(newlist, v)
		}
	}
	return newlist
}

// check if host is active ?
func Ishostactive(host string) bool {
	client, err := goph.NewUnknown("root", host, goph.Password(Getpass()))
	if err != nil {
		return false
	}
	client.Close()
	return true
}

// writedata updates/rewrites data into file
func Update(file string, list []string) error {
	data := ""
	for _, item := range list {
		data += item + "\n"
	}
	err := os.WriteFile(file, []byte(data), 0644)
	if err != nil {
		log.Println(err)
	}
	return err
}

// filterlist make list unique
func Unique(list []string) []string {
	mp := make(map[string]bool)
	for _, h := range list {
		mp[h] = true
	}
	ulist := make([]string, 0)
	for h := range mp {
		if h == "" {
			break
		}
		ulist = append(ulist, h)
	}
	return ulist
}

// checkerr check error if exeste and Close program
func checkerr(info string, err error) {
	if err != nil {
		log.Println(info, err)
		os.Exit(0)
	}
}
