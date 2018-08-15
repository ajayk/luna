package customio

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

func ReplaceToml(filepath, moniker, externalip string) {
	exAddrStr := `external_address = "` + externalip + `:26656"`
	monikerStr := `moniker = "` + moniker + `"`
	input, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Index(line, "external_address") == 0 {
			lines[i] = exAddrStr
		} else if strings.Index(line, "moniker") == 0 {
			lines[i] = monikerStr
		}
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filepath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

// Return ip address of 1) seeds 2) persistent_peers 3) private_peers_ids from config.toml
func ReadToml(filepath string, target *[]string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		readToml(scanner.Text(), target)
	}
}

func readToml(msg string, target *[]string) {
	acceptIpsDns := strings.Split(msg, "\n")
	for _, line := range acceptIpsDns {
		if line != "" {
			runes := []rune(line)
			if string(runes[0:5]) == "seeds" || string(runes[0:16]) == "persistent_peers" {
				linesplit := strings.Split(line, "=")
				peersstr := linesplit[1]
				strings.Replace(peersstr, `"`, "", -1) // remove " from string
				trimpeersStr := strings.Trim(peersstr, " ")

				if !strings.Contains(trimpeersStr, "@") {
					continue
				}

				idandips := strings.Split(trimpeersStr, ",")
				for _, idandip := range idandips {
					idandipsplit := strings.Split(idandip, "@")
					ipandport := idandipsplit[1]
					ipandportsplit := strings.Split(ipandport, ":")
					ipordns := ipandportsplit[0]
					ip := net.ParseIP(ipordns)
					if ip != nil {
						//log.Println("step1", ip)
						*target = append(*target, ip.String())
					} else {
						ips, err := net.LookupIP(ipordns)
						if err != nil {
							log.Println(os.Stderr, "Could not get IPs::", ipordns, err)
						}
						for _, ip := range ips {
							//log.Println("step2", ip)
							*target = append(*target, ip.String())
						}
					}
				}
			}
		}
	}
}
