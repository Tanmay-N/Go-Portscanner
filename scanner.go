package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/tatsushid/go-fastping"
)

type response struct {
	addr *net.IPAddr
	rtt  time.Duration
}

func main() {
	textPtr := flag.String("IP", "127.0.0.1", "IP Address/Domain name. (Required)")
	flag.Parse()

	if *textPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	FastScan(*textPtr)
}

//FastScan takes IP as a String
func FastScan(ip string) {
	var ipList []string

	if strings.Contains(ip, "-") {
		ipList = append(ipList, ParseIPSequence(ip)...)
	} else {
		ipList = []string{ip}
	}

	for _, i := range ipList {

		if osChek(i) {
			fmt.Println("scanning started", i)
			scanresult := `----------------------------
        SCAN RESULTS
----------------------------`
			fmt.Println(scanresult)
			target := ip
			start := time.Now()
			activeThreads := 0
			doneChannel := make(chan bool)
			fmt.Println("Port", "\tStatus", "\tService\t\t")
			fmt.Println("----", "\t------", "\t------\t\t")
			val := 0

			for port := 0; port <= 65535; port++ {
				go scanTCPConnection(target, port, doneChannel, &val)
				time.Sleep(time.Millisecond)
				activeThreads++
			}

			// Wait for all threads to finish
			for activeThreads > 0 {
				<-doneChannel
				activeThreads--
			}
			fmt.Println("\n---------------------------")
			fmt.Println(val, "Ports are Open")
			fmt.Println("Time Took: ", time.Since(start))

		} else {
			fmt.Println("Not able to Reach the Domain/IP", i)
			//os.Exit(1)
		}
	}

}

func scanTCPConnection(ip string, port int, doneChannel chan bool, val *int) {
	_, err := net.DialTimeout("tcp", ip+":"+strconv.Itoa(port),
		time.Second*1)
	if err == nil {
		fmt.Println(port, "\tOpen\t", portShortList[ToString(port)])
		*val = *val + 1
	}
	doneChannel <- true
}

//Ping the Domain and return Bool Value (For Windows)
func resolveIP(host string) bool {

	hostname := host
	if len(hostname) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	p := fastping.NewPinger()
	netProto := "ip4:icmp"
	if strings.Index(hostname, ":") != -1 {
		netProto = "ip6:ipv6-icmp"
	}
	ra, err := net.ResolveIPAddr(netProto, hostname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	results := make(map[string]*response)
	results[ra.String()] = nil
	p.AddIPAddr(ra)

	onRecv, onIdle := make(chan *response), make(chan bool)
	p.OnRecv = func(addr *net.IPAddr, t time.Duration) {
		onRecv <- &response{addr: addr, rtt: t}
	}
	p.OnIdle = func() {
		onIdle <- true
	}

	p.MaxRTT = time.Second
	p.RunLoop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

loop:
	for {
		select {
		case <-c:
			fmt.Println("get interrupted")
			break loop
		case res := <-onRecv:
			if _, ok := results[res.addr.String()]; ok {
				results[res.addr.String()] = res
			}
		case <-onIdle:
			for _, r := range results {
				if r == nil {
					return false
				}
				return true
			}
		case <-p.Done():
			if err = p.Err(); err != nil {
				fmt.Println("Ping failed:", err)
			}
			break loop
		}
	}
	signal.Stop(c)
	p.Stop()
	return false
}

//ParseIPSequence is Ip scanning
func ParseIPSequence(ipSequence string) []string {

	var arrayIps []string

	series, _ := regexp.Compile("([0-9]+)")

	// For sequence ips, using '-'
	lSeries := series.FindAllStringSubmatch(ipSequence, -1)

	for i := ToInt(lSeries[3][0]); i <= ToInt(lSeries[4][0]); i++ {
		arrayIps = append(arrayIps,
			lSeries[0][0]+"."+
				lSeries[1][0]+"."+
				lSeries[2][0]+"."+ToString(i))
	}

	//fmt.Println(lSeries[3][0])
	//	fmt.Println(lSeries[4][0])
	return arrayIps
}

//ToInt is to convert into integer
func ToInt(s string) int {

	i, _ := strconv.Atoi(s)
	return i
}

//ToString is to convert into String
func ToString(s int) string {

	i := strconv.Itoa(s)
	return i
}

func osChek(ip string) bool {
	switch os := runtime.GOOS; os {
	case "linux":
		return pingIP(ip)
	case "darwin":
		return pingIP(ip)
	default:
		fmt.Println("Not supported for the OS")
		return resolveIP(ip)
	}

}

func pingIP(ip string) bool {
	Command := fmt.Sprintf("ping -c 1 " + ip + " > /dev/null && echo true || echo false")
	output, err := exec.Command("/bin/sh", "-c", Command).Output()
	test := string(output)
	if strings.Contains(test, "true") {
		return true
	}

	if err != nil {
		fmt.Println("found error")
		return false
	}
	return false

}
