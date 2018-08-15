package install

import (
	"log"

	"github.com/lunamint/luna/util/shell"
)

func IfstatInstall() {
	ifstat := "sudo apt install ifstat -y"
	msg, err := shell.SimpleShellCall("", ifstat, false)
	if err != nil {
		log.Fatal("ifstat install  error: ", msg, err)
	} else {
		log.Println("ifstat install done...")
	}
}
