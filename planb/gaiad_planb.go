package planb

import (
	"log"

	"strconv"
	"strings"
	"time"

	"github.com/lunamint/luna/util/awsobject"
	"github.com/lunamint/luna/util/customio"
	"github.com/lunamint/luna/util/install"
	"github.com/lunamint/luna/util/shell"
)

func assginNewElsaticIP(region, instanceID string) string {
	// release existing elastic ips*allocRes.PublicIp
	awsobject.ReleaseElasticip(region)
	return awsobject.AllocateIP(region, instanceID)
}

func gaiadStart() {
	startcmd := "$HOME/bin/gaiad_start.sh"
	msg, err := shell.SimpleShellCall("", startcmd, false)
	if err != nil {
		log.Println("gaiad start error:", msg, err)
	} else {
		log.Println("gaiad started")
	}
}

func gaiadStop() {
	stopcmd := "$HOME/bin/gaiad_stop.sh"
	msg, err := shell.SimpleShellCall("", stopcmd, false)
	if err != nil {
		log.Println("gaiad stop  error:", msg, err)
	} else {
		log.Println("gaiad stopped")
	}
}

func getNetworkInterface() string {
	interfacecmd := "route | grep '^default' | grep -o '[^ ]*$'"
	msg, err := shell.SimpleShellCall("", interfacecmd, false)
	if err != nil {
		log.Println("get newtork interface error:", msg, err)
	}
	return msg
}

func Bplan(home, region, gaiad_config_toml, planb_topic_arn string) {
	install.SysstatInstall()
	install.IfstatInstall()

	bPlanFor(home, region, gaiad_config_toml, planb_topic_arn)
}

func bPlanFor(home, region, gaiad_config_toml, planb_topic_arn string) {
	// consecutive 3 strikes out
	var cpustrikes int = 0
	var netstrikes int = 0

	//backup files
	pathBak := home + "/.gaiad/config/"
	pathTmp := home
	addrbook := "addrbook.json"
	config := "config.toml"
	genesis := "genesis.json"
	var backups []string
	backups = append(backups, addrbook)
	backups = append(backups, config)
	backups = append(backups, genesis)

	// find out networkt interface
	netinter := getNetworkInterface()

	for {
		// network usage
		//netusgcmd := "b(){ echo $(($(cat /sys/class/net/" + netinter + "/statistics/tx_bytes)+$(cat /sys/class/net/" + netinter + "/statistics/rx_bytes))); }; sn(){ sleep 0.$((1000000000-$(date '+%N'|sed 's/0*//'))); }; lb=$(b); nb=$(b); echo $((($nb-$lb)/1024));lb=$nb; sn; done"
		netusgcmd := "ifstat -i " + netinter + " -q 1 1"
		netusgmsg, err := shell.SimpleShellCall("", netusgcmd, false)
		if err != nil {
			log.Println("network usage calculation error:", netusgmsg, err)
		}

		netmsglines := strings.Split(netusgmsg, "\n")
		inoutusg := netmsglines[2]

		inouts := strings.Split(inoutusg, " ")
		var input float64 = -1
		var output float64 = -1

		for _, inout := range inouts {
			if inout != "" {
				if input == -1 {
					input, err = strconv.ParseFloat(inout, 64)
					if err != nil {
						log.Fatal("input conversion error: ", err)
					}
				} else {
					output, err = strconv.ParseFloat(inout, 64)
					if err != nil {
						log.Fatal("output conversion error: ", err)
					}
				}
			}
		}

		// cpu usage check
		mpcmd := "mpstat | awk '$3 ~ /CPU/ { for(i=1;i<=NF;i++) { if ($i ~ /%idle/) field=i } } $3 ~ /all/ { print 100 - $field }'"
		cpumsg, err := shell.SimpleShellCall("", mpcmd, false)
		if err != nil {
			log.Println("cpu usage calculation error:", cpumsg, err)
		}

		cpuUsage, err := strconv.ParseFloat(cpumsg, 64)
		if err != nil {
			log.Println("strconv at serverWatcher error: ", err)
		}

		if cpuUsage > 70 { //percent
			cpustrikes++
		} else {
			if cpustrikes > 0 {
				cpustrikes--
			}
		}

		if input > 1500000 || output > 1500000 { // KB, euqal to 1.5 GB
			netstrikes++
		} else {
			if netstrikes > 0 {
				netstrikes--
			}
		}

		//		log.Println("input=", input)
		//		log.Println("output=", output)
		//		log.Println("cpuUsage=", cpuUsage)
		//		log.Println("cpustrikes=", cpustrikes)
		//		log.Println("netstrikes=", netstrikes)

		if cpustrikes >= 3 || netstrikes >= 3 { // cpu usage is bigger then 70  occured 3 consecutive times or network usage is bigger then 8GB 3 consecutive times
			log.Println("Cpu/Bandwidth usage not normal.. do planB , gaiad init & restart")

			nodeidcmd := "gaiad tendermint show_node_id"
			oldnodeid, err := shell.SimpleShellCall("", nodeidcmd, false)
			if err != nil {
				log.Println("gaiad tendermint show_node_id error:", err)
			}

			myipcmd := "curl ifconfig.co"
			oldip, err := shell.SimpleShellCall("", myipcmd, false)
			if err != nil {
				log.Println("curl ifconfig.co error:", err)
			}

			oldAddr := oldnodeid + "@" + oldip

			// gaiad stop
			gaiadStop()

			// backup files before init
			customio.Backupfiles(pathBak, backups, pathTmp)

			// delete gaiad config directory
			customio.DeleteConfigDirectory(pathBak)

			// get random byte for node name
			ranint := customio.RandomInt(4, 9)
			moniker := string(customio.RandASCIIBytes(ranint))

			// gaiad  init
			//gaiadcmd := "gaiad init --name=" + string(moniker) + " --chain-id=" + chainid
			gaiadcmd := "gaiad init --name=" + moniker
			msg, err := shell.SimpleShellCall("", gaiadcmd, false)
			if err != nil {
				log.Println("gaiad init  error:", msg, err)
			} else {
				log.Println(gaiadcmd)
				log.Println("gaiad initiated: ", msg)
			}

			// wrtie secret msg to file
			customio.WriteToFile(home+"/secret", msg)

			//restore files
			customio.Restorefiles(pathBak, backups, pathTmp)

			// assign new elasticip
			eip := assginNewElsaticIP(region, awsobject.GetInstanceID(region))

			// update config.toml's external ip as newelasticip
			customio.ReplaceToml(gaiad_config_toml, moniker, eip)

			// gaiad start
			gaiadStart()

			// new nodeid
			newnodeid, err := shell.SimpleShellCall("", nodeidcmd, false)
			if err != nil {
				log.Println("gaiad tendermint show_node_id error:", err)
			}

			newAddr := newnodeid + "@" + eip

			//publish message
			PublishMessage(region, "sentry ip changed", oldAddr+"\n"+newAddr, planb_topic_arn)

			cpustrikes = 0
			netstrikes = 0
		}

		//sleep
		time.Sleep(10 * time.Second)
	}
}
