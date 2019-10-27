package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

//IPv4 asdasda
type IPv4 [4]int

func main() {
	var portsList []string
	var IPsListapp []string

	var ports = flag.String("port", "22", "Specify the port")
	var ips = flag.String("ip", "scanme.nmap.org", "Specify the ipaddress")
	all := flag.Bool("A", false, "Scans from port 1 to 1024")
	flag.Parse()

	if *all {
		*ports = "1-1024"
	}

	portsList = append(portsList, *ports)
	IPsListapp = append(IPsListapp, *ips)

	IPScanner(IPsListapp, portsList, true)

}

//IPScanner asda
func IPScanner(ipstr []string, portStr []string, printResults bool) map[IPv4][]string {

	m := make(map[IPv4][]string)

	var ipList []IPv4
	var portList []string

	var wg sync.WaitGroup
	//regMail, _ := regexp.Compile(`[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,6}`)
	if len(portStr) == 1 {
		portList = ParsePortList(portStr[0])
	} else {
		portList = portStr
	}

	if len(ipstr) == 0 {
		ipList = append(ipList, IPv4{127, 0, 0, 1})
	} else {
		for _, i := range ipstr {
			if strings.Contains(i, "-") {
				ipList = append(ipList, ParseIPSequence(i)...)
			} else {
				ip := ToIPv4(i)
				if ip.IsValid() {
					ipList = append(ipList, ip)
				} else {
					ipList = append(ipList, domainToIP(i))
				}
			}
		}
	}

	for _, ip := range ipList {
		wg.Add(1)
		go func(ip IPv4) {
			defer wg.Done()
			result := PortScanner(ip, portList)
			if len(result) > 0 {
				m[ip] = result
				if printResults {
					PresentResults(ip, result)
				}
			} else {
				fmt.Println("nothing found", result)
			}
		}(ip)
	}

	wg.Wait()

	return m
}

//PortScanner is scanning
func PortScanner(ip IPv4, portList []string) []string {

	var open []string

	for _, port := range portList {

		_, err := net.DialTimeout("tcp", ip.ToString()+":"+port, time.Second*50)

		if err == nil {
			//conn.Close()
			open = append(open, port)
		}
	}

	return open
}

//Check checking error
func Check(err error) {

	if err != nil {
		panic(err)
	}
}

// ToInt converts a string to integer, as strconv.Atoi does, but without
// returning errors.
func ToInt(s string) int {

	i, _ := strconv.Atoi(s)
	return i
}

// ToString converts an IP from IPv4 type to string.
func (ip *IPv4) ToString() string {

	ipStringed := strconv.Itoa(ip[0])
	for i := 1; i < 4; i++ {
		strI := strconv.Itoa(ip[i])
		ipStringed += "." + strI
	}
	return ipStringed
}

// IsValid checks an IP address as valid or not.
func (ip *IPv4) IsValid() bool {

	for i, oct := range ip {
		if i == 0 || i == 3 {
			if oct < 1 || oct > 254 {
				return false
			}
		} else {
			if oct < 0 || oct > 255 {
				return false
			}
		}
	}
	return true
}

// PlusPlus increments an IPv4 value.
func (ip *IPv4) PlusPlus() *IPv4 {

	if ip[3] < 254 {
		ip[3] = ip[3] + 1
	} else {
		if ip[2] < 255 {
			ip[2] = ip[2] + 1
			ip[3] = 1
		} else {
			if ip[1] < 255 {
				ip[1] = ip[1] + 1
				ip[2] = 1
				ip[3] = 1
			} else {
				if ip[0] < 255 {
					ip[0] = ip[0] + 1
					ip[1] = 1
					ip[2] = 1
					ip[3] = 1
				}
			}
		}
	}
	return ip
}

// ToIPv4 converts an string to a IPv4. ans [8 8 8 8]
//https://play.golang.org/p/kvAPyjkta1f
func ToIPv4(ip string) IPv4 {

	var newIP IPv4

	ipS := strings.Split(ip, ".")

	for i, v := range ipS {
		newIP[i], _ = strconv.Atoi(v)
	}

	return newIP
}

// ParseIPSequence gets a sequence of IP addresses correspondent from an
// "init-end" entry.
func ParseIPSequence(ipSequence string) []IPv4 {

	var arrayIps []IPv4

	series, _ := regexp.Compile("([0-9]+)")

	// For sequence ips, using '-'
	lSeries := series.FindAllStringSubmatch(ipSequence, -1)

	for i := ToInt(lSeries[3][0]); i <= ToInt(lSeries[4][0]); i++ {
		arrayIps = append(arrayIps, IPv4{
			ToInt(lSeries[0][0]),
			ToInt(lSeries[1][0]),
			ToInt(lSeries[2][0]),
			i})
	}
	return arrayIps
}

// ParsePortList gets a port list from its port entry in arguments.
func ParsePortList(rawPorts string) []string {

	var ports []string

	individuals, _ := regexp.Compile("([0-9]+)[,]*")
	series, _ := regexp.Compile("([0-9]+)[-]([0-9]+)")

	// For individual ports, separated by ','
	lIndividuals := individuals.FindAllStringSubmatch(rawPorts, -1)

	// For sequence ports, using '-'
	lSeries := series.FindAllStringSubmatch(rawPorts, -1)

	if len(lSeries) > 0 {
		for _, s := range lSeries {
			init, _ := strconv.Atoi(s[1])
			end, _ := strconv.Atoi(s[2])
			for i := init + 1; i < end; i++ {
				ports = append(ports, strconv.Itoa(i))
			}
		}
	}
	for _, port := range lIndividuals {
		ports = append(ports, port[1])
	}
	sort.Strings(ports)

	return ports
}

// GetAllIPsClassC returns a slice of IPv4 with all IP addresses
// from a Class C.
//func GetAllIPsClassC(ip IPv4) []IPv4 {}

// PresentResults presents all results in console.
func PresentResults(ip IPv4, ports []string) int {

	fmt.Println(" \n>" + ip.ToString())
	fmt.Println(" Port:	Description:")
	for _, port := range ports {
		fmt.Println(" " + port + "\t" + portShortList[port])
	}
	return 0
}

func domainToIP(ip string) IPv4 {
	var arrayIps IPv4

	ip4address, err := net.ResolveIPAddr("ip4", ip)
	if err != nil {
		fmt.Println("Fail to resolve IP4", err.Error())
		os.Exit(1)
	}

	i := ip4address.String()

	arrayIps = ToIPv4(i)

	return arrayIps

}
