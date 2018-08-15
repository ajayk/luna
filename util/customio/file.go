package customio

import (
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/lunamint/luna/util/shell"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandASCIIBytes(n int) []byte {
	output := make([]byte, n)
	// We will take n bytes, one byte for each character of output.
	randomness := make([]byte, n)
	// read all random
	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}
	l := len(letterBytes)
	// fill output
	for pos := range output {
		// get random item
		random := uint8(randomness[pos])
		// random % 64
		randomPos := random % uint8(l)
		// put into output
		output[pos] = letterBytes[randomPos]
	}
	return output
}

func Deletefiles(path string, files []string) {
	for _, file := range files {
		sudodelete(file)
	}
}

func DeleteConfigDirectory(path string) {
	deletecmd := "sudo rm -rf " + path
	msg, err := shell.SimpleShellCall("", deletecmd, false)
	if err != nil {
		log.Println("sudodelete error:", msg, err)
	}
}

func sudodelete(obj string) {
	deletecmd := "sudo rm -rf " + obj
	msg, err := shell.SimpleShellCall("", deletecmd, false)
	if err != nil {
		log.Println("sudodelete error:", msg, err)
	}
}

func Backupfiles(path string, files []string, pathTmp string) {
	for _, file := range files {
		source := path + file
		dest := pathTmp + "/" + file + ".bak"
		backupcmd := "sudo cp " + source + " " + dest
		msg, err := shell.SimpleShellCall("", backupcmd, false)
		if err != nil {
			log.Println("backupcmd error: ", msg, err)
		}
		sudodelete(source)
	}
}

func Restorefiles(path string, files []string, pathTmp string) {
	for _, file := range files {
		source := pathTmp + "/" + file + ".bak"
		dest := path + file
		restorecmd := "sudo cp " + source + " " + dest
		msg, err := shell.SimpleShellCall("", restorecmd, false)
		if err != nil {
			log.Println("restorecmd error: ", msg, err)
		}
		sudodelete(source)
	}
}

func WriteToFile(path string, content string) (err error) {
	var f *os.File
	f, err = os.Create(path)
	if err != nil {
		return
	}

	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()

	_, err = io.WriteString(f, content)
	return
}

func RandomInt(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
