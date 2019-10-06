package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"regexp"
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
	textPtr := flag.String("IP", "", "IP Address/Domain name. (Required)")
	flag.Parse()

	if *textPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fastScan(*textPtr)
}

func fastScan(ip string) {
	regMail, _ := regexp.Compile(`[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,6}`)
	if resolveIP(ip) || regMail.MatchString(string(ip)) {
		scanresult := `----------------------------
	SCAN RESULTS
----------------------------`
		fmt.Println(scanresult)
		//	time.Sleep(time.Millisecond * 10)
		target := ip
		start := time.Now()
		activeThreads := 0
		doneChannel := make(chan bool)
		fmt.Println("Port", "\tStatus\t\t")
		fmt.Println("----", "\t------\t\t")
		val := 0
		for port := 0; port <= 65535; port++ {
			go scanTCPConnection(target, port, doneChannel, &val)
			time.Sleep(1 * time.Nanosecond)
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
		fmt.Println("Not able to Reach the Domain/IP", ip)
		os.Exit(1)
	}

}

func scanTCPConnection(ip string, port int, doneChannel chan bool, val *int) {
	_, err := net.DialTimeout("tcp", ip+":"+strconv.Itoa(port),
		time.Second*10)
	if err == nil {
		fmt.Println(port, "\tOpen\t\t")
		*val = *val + 1
	}
	doneChannel <- true
}

//Ping the Domain and return Bool Value
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
