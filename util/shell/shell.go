package shell

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"
)

var (
	// stepDepth         = 0
	separatorLogGroup = strings.Repeat("=", 60)
	separatorPanic    = strings.Repeat("=", 80)
)

func RealtimeShellCall(cmdStr string) {
	cmd := exec.Command(os.Getenv("SHELL"), "-c", cmdStr)
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	oneRune := make([]byte, utf8.UTFMax)
	for {
		count, err := stdout.Read(oneRune)
		if err != nil {
			break
		}
		fmt.Printf("%s", oneRune[:count])
	}
}

// SimpleShellCall info
//==============================================================================
func SimpleShellCall(workDir string, cmdStr string, isUsePanic bool) (msg string, err error) {
	var stderr bytes.Buffer
	c := exec.Command(os.Getenv("SHELL"), "-c", cmdStr)
	c.Stderr = &stderr

	if workDir != "" {
		c.Dir = workDir
	}

	cmdOut, err := c.Output()

	if err == nil {
		msg = strings.TrimSpace(string(cmdOut[:]))
		// fmt.Printf(">>> %s\n", msg)
		return
	}

	if isUsePanic {
		msg = "[ERR] " + err.Error() + "\n" +
			"[DIR] " + workDir + "\n" +
			"[CMD] " + cmdStr + "\n" +
			"[MSG] " + string(stderr.Bytes())
		MakePanic("Error calling shell command", msg)
	}

	msg = string(stderr.Bytes())
	return
}

// MakePanic info
//==============================================================================
func MakePanic(title string, contents string) {
	log.Println(separatorLogGroup)
	e := "<< " + title + " >>\n" + separatorPanic + "\n" + contents + "\n" + separatorPanic
	panic(e)
}
