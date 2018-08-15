package main

import (
	//"fmt"
	"log"
	"os"

	"github.com/lunamint/luna/util/awsobject"
	"github.com/lunamint/luna/util/install"

	"github.com/lunamint/luna/net"
	"github.com/lunamint/luna/planb"
)

var (
	HOME              string
	PLANB_TOPIC_ARN   string
	GAIAD_ADDRBOOK    string
	GAIAD_CONFIG_TOML string
	GAIAD_GENESIS     string
	AWS_REGION        string
)

func init() {
	// AWS_REGION
	AWS_REGION = awsobject.GetRegion()

	//HOME
	if os.Getenv("HOME") == "" {
		HOME = "/home/ubuntu"
	} else {
		HOME = os.Getenv("HOME")
	}

	// PLANB_TOPIC_ARN
	PLANB_TOPIC_ARN = "arn:aws:sns:" + AWS_REGION + ":" + awsobject.GetAccountID(AWS_REGION) + ":sipchanged"

	GAIAD_ADDRBOOK = HOME + "/.gaiad/config/addrbook.json"
	GAIAD_CONFIG_TOML = HOME + "/.gaiad/config/config.toml"
	GAIAD_GENESIS = HOME + "/.gaiad/config/genesis.json"
}

func isAWSCredentialExist() bool {
	if _, err := os.Stat(HOME + "/.aws"); err == nil {
		return true
	} else {
		return false
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("needs args [cmds] [planb] [net]")
	}

	if !isAWSCredentialExist() {
		log.Fatal("AWS credentail not exist at " + HOME + "./aws")
	}

	//Make Bin, gaiad_start.sh, gaiad_stop.sh
	install.MakeDirectories(HOME)

	cmd := os.Args[1]

	if cmd == "cmds" {
		log.Println("luna planb &> $HOME/bin/planb & ")
		return
	}

	if cmd == "planb" {
		log.Println("start planb..")
		planb.Bplan(HOME, AWS_REGION, GAIAD_CONFIG_TOML, PLANB_TOPIC_ARN)
		return
	}

	if cmd == "net" {
		log.Println("start netstat geo status tracking..")
		net.NetStatGeoStatus(AWS_REGION)
		return
	}

	log.Println("wrong cmd!")
}
