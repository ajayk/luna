package net

import (
	"log"
	"strings"

	"github.com/lunamint/luna/util/awsobject"
	"github.com/lunamint/luna/util/color"
	"github.com/lunamint/luna/util/install"
	"github.com/lunamint/luna/util/shell"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func getSentryIplist(svc *ec2.EC2, target *[]string) {
	for _, g := range awsobject.DescribeSecurityGroups(svc).SecurityGroups {
		if *g.GroupName != "SSH" {
			for i := 0; i < len(g.IpPermissions); i++ {
				gper := g.IpPermissions[i]
				for j := 0; j < len(gper.IpRanges); j++ {
					iprange := gper.IpRanges[j]
					// remove "/32"
					pureip := strings.Replace(*iprange.CidrIp, "/32", "", -1)
					//log.Println(pureip)
					if pureip != "0.0.0.0/0" {
						*target = append(*target, pureip)
					}
				}
			}
		}
	}
}

func NetStatGeoStatus(aws_region string) {
	if !install.GeolookupInstalled() {
		install.GeolookupInstall()
	}
	printStatus(aws_region)
}

func printStatus(aws_region string) {
	cmd := "netstat -ntp | grep gaiad | grep ESTABLISHED"
	msg, err := shell.SimpleShellCall("", cmd, false)
	if err != nil {
		log.Println("status cmd error:", err)
	}

	var privateip string
	curlcmd := "curl http://169.254.169.254/latest/meta-data/local-ipv4"
	curlmsg, err := shell.SimpleShellCall("", curlcmd, false)
	if err != nil {
		log.Println("curlcmd error:", err)
	} else {
		privateip = curlmsg
	}

	//log.Println("privateip=", privateip)
	var mAddrList map[string]bool
	mAddrList = make(map[string]bool)

	acceptIpsDns := strings.Split(msg, "\n")
	for _, line := range acceptIpsDns {
		//log.Println("line=", line)
		if line != "" {
			words := strings.Split(line, " ")
			for _, word := range words {
				if strings.Contains(word, ":") {
					swords := strings.Split(word, ":")
					if swords[0] != privateip {
						if mAddrList[swords[0]] == false {
							mAddrList[swords[0]] = true
						}
					}
				}
			}
		}
	}

	for ip, _ := range mAddrList {
		geocmd := "geoiplookup -f /usr/local/share/GeoIP/GeoLiteCity.dat " + ip
		msg, err := shell.SimpleShellCall("", geocmd, false)
		if err != nil {
			log.Println("geolookup error:", err)
		} else {
			// log.Println(ip, "	", msg)
			log.Printf("%-15s %s", ip, msg)
		}
	}

	color.Blue("-----------------------")
	log.Println("total:", len(mAddrList))
}
