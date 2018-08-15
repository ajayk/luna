package install

import (
	"log"
	"os"

	"github.com/lunamint/luna/util/shell"
)

func SysstatInstall() {
	sysstat := "sudo apt install  sysstat -y"
	msg, err := shell.SimpleShellCall("", sysstat, false)
	if err != nil {
		log.Fatal("sysstat install  error: ", msg, err)
	} else {
		log.Println("sysstat install done...")
	}
}

func GeolookupInstalled() bool {
	if _, err := os.Stat("/usr/local/share/GeoIP/GeoLiteCity.dat"); err == nil {
		return true
	} else {
		return false
	}
}

func GeolookupInstall() {
	geo := "sudo apt install  geoip-bin -y"
	msg, err := shell.SimpleShellCall("", geo, false)
	if err != nil {
		log.Println("geoiplookup install  error: ", msg, err)
	}

	getcitydat := "sudo wget -N http://geolite.maxmind.com/download/geoip/database/GeoLiteCity.dat.gz"
	msg, err = shell.SimpleShellCall("", getcitydat, false)
	if err != nil {
		log.Println("wget  geolitecity.dat error: ", msg, err)
	}

	gzip := "sudo gunzip GeoLiteCity.dat.gz"
	msg, err = shell.SimpleShellCall("", gzip, false)
	if err != nil {
		log.Println("gunzip  geolitecity.dat  error: ", msg, err)
	}

	mkdir := "sudo mkdir /usr/local/share/GeoIP/"
	msg, err = shell.SimpleShellCall("", mkdir, false)
	if err != nil {
		log.Println("mkdir /usr/local/shre/Geoip error: ", msg, err)
	}

	mvdat := "sudo mv ./GeoLiteCity.dat /usr/local/share/GeoIP/"
	msg, err = shell.SimpleShellCall("", mvdat, false)
	if err != nil {
		log.Println("move geolitecity.dat error: ", msg, err)
	}
}
