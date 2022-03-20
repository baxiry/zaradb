package helper

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/melbahja/goph"
)

//  unzipRemote unzip remote zipped file
func UnzipRemote(sshclient *goph.Client, zippedfile string) error {
	// zip the client bot app
	cmd, err := sshclient.Command("unzip", "-o", "/root/"+zippedfile+"-bot.zip")
	//cmd, err := sshclient.Command("tar", "-xf", "/root/"+zippedfile+"-bot.zip")
	if err != nil {
		return err
	}
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// TODO test localzip function
//  zipfile.zip and clientName
func ZipLocalDir(source string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(source + "-bot.zip")
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source+"-bot", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		if err != nil {
			return err
		}
		return nil

	})
}

// deploy deploy client-bot.zip to client host
func Deploy(clientBot, hostBot string) error {
	sshClient, err := goph.NewUnknown("root", hostBot, goph.Password(Getpass()))
	if err != nil {
		return err
	}

	clientBot = clientBot + "-bot.zip"

	fmt.Println(clientBot)
	err = sshClient.Upload(clientBot, clientBot)
	if err != nil {
		return err
	}

	return nil
}

// to scure app read pass form seprite file
func Getpass() string {
	data, err := ioutil.ReadFile(".mypass")
	if err != nil {
		return err.Error()
	}
	psw := string(data)
	return psw[:len(psw)-1]
}

// todo context pkg must be used in this function
// runrmotebot runc remote bot app
// note that botdir same clinet name
func Runrmotebot(host, botdir string) error {
	sshclient, err := goph.NewUnknown("root", host, goph.Password(Getpass()))
	if err != nil {
		fmt.Println("err when connete")
		return err
	}
	defer sshclient.Close()

	_, err = sshclient.Run("/root/" + botdir + "-bot/testbot &")
	if err != nil {
		fmt.Println("err when ")
		return err
	}
	return nil
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
func Copylocaldir(src string, dst string) error {
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
			if err = Copylocaldir(srcfp, dstfp); err != nil {

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

// run
// todo test zip function
//  zipfile.zip and clientname
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

// load loads file and return hosts address as []stirng
func Load(file string) ([]string, error) {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	hs := strings.Split(string(data), "\n")

	list := make([]string, 0)

	for _, v := range hs {

		h := strings.Replace(v, " ", "", -1)
		if len(h) < 3 {
			continue
		}
		list = append(list, h)
	}

	return Unique(list), nil
}

// appendaddress appends new address to addressfile
func Appendaddr(file, data string) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(data + "\n"); err != nil {
		log.Println(err)
	}
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

// exitbot
func Exitbot() {
	ch := make(chan bool, 1)
	http.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "will done") //": %s\n", r.url.path)
		ch <- true
	})
	go func() {
		fmt.Println(http.ListenAndServe(":80", nil))
	}()

	go func() {
		if <-ch {
			fmt.Println("done")
			os.Exit(0)
		}
	}()
}

// send exitbot
func Sendexit(address string) {
	resp, err := http.Get("http://" + address + "/exit")
	if err != nil {
		log.Fatal("error getting response. ", err)
	}
	defer resp.Body.Close()

	// read body from response
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("error reading response. ", err)
	}

	fmt.Printf("body is : %s\n", body)
}

// may be not need this func
// copies a file. and rename to name with .cp saffix
func Copyfile(src string) error {
	// open the source file for reading
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	// open the destination file for writing
	d, err := os.Create(src + ".cp")
	if err != nil {
		return err
	}

	// Copy the contents of the source file into the destination file
	if _, err := io.Copy(d, source); err != nil {
		d.Close()
		return err
	}

	// return any errors that result from closing the destination file
	// will return nil if no errors occurred
	return d.Close()
}

// checkerr check error if exeste and Close program
func checkerr(info string, err error) {
	if err != nil {
		log.Println(info, err)
		os.Exit(0)
	}
}
