package install

import (
	"log"
	"os"

	"github.com/lunamint/luna/util/customio"
	"github.com/lunamint/luna/util/shell"
)

func MakeDirectories(home string) error {
	bin := home + "/bin"
	gaialog := home + "/.gaiad_log"
	//check bin directory ,create if not exist
	if _, err := os.Stat(bin); err != nil {
		mkdir := "mkdir " + bin
		msg, err := shell.SimpleShellCall("", mkdir, false)
		if err != nil {
			log.Println("mkdir "+bin+" error: ", msg, err)
			return err
		}
	}

	//check .gaiad_log directory, create if not exist
	if _, err := os.Stat(gaialog); err != nil {
		mkdir := "mkdir " + gaialog
		msg, err := shell.SimpleShellCall("", mkdir, false)
		if err != nil {
			log.Println("mkdir "+gaialog+" error: ", msg, err)
			return err
		}
	}

	//create if not exist
	if _, err := os.Stat(bin + "/gaiad_start.sh"); err != nil {
		createFile(bin, "gaiad_start.sh", "744", getStartStr())
	}

	//create if not exist
	if _, err := os.Stat(bin + "/gaiad_stop.sh"); err != nil {
		createFile(bin, "gaiad_stop.sh", "744", getStopStr())
	}

	return nil
}

func createFile(path string, file string, auth string, content string) {
	err := customio.WriteToFile(path+"/"+file, content)
	if err != nil {
		shell.MakePanic("write file error", path+"/"+file)
	}
	shell.SimpleShellCall(path, "chmod "+auth+" "+file, true)
}

func getStartStr() string {
	return `#!/usr/bin/env bash

is_run=` + "`" + `ps -ef | grep "gaiad start" | grep -v grep | wc -l` + "`" + `

if [ ${is_run} -ne 0 ]; then
        echo "gaiad is already running"
        exit 1
fi

gaiad start >> $HOME/.gaiad_log/gaiad.log 2>&1 &
echo "Done"
`
}

func getStopStr() string {
	return `#!/usr/bin/env bash

is_run=` + "`" + `ps -ef | grep "gaiad start" | grep -v grep | wc -l` + "`" + `

if [ ${is_run} -eq 0 ]; then
        echo "gaiad is not running!!"
        exit 1
fi

ps -ef | grep "gaiad start" | grep -v grep | awk '{print $2}' | xargs kill -9
echo "Done"
`
}
