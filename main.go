package main

import (
	"os"
	"fmt"
	"net"
	"time"
	"sync"
	"strings"
	"strconv"
)

type Scan struct {
	host string
	port string
	status int
}

type Results struct {
	mutex sync.Mutex
	scans []Scan
}

func parsePorts(portsRaw string) []string {

	var ports []string

	if (strings.Contains(portsRaw, ":")) {
		// range of ports
		portSplit := strings.Split(portsRaw, ":")
		firstPort, _ := strconv.Atoi(portSplit[0])
		secondPort, _ := strconv.Atoi(portSplit[1])
		for i := firstPort; i <= secondPort; i++ {
			ports = append(ports, strconv.Itoa(i))
		}

	} else {
		// single port
		ports = append(ports, portsRaw)
	}

	return ports
}

func scanPort(host string, port string, results *Results) {
	var status int

	conn, err := net.DialTimeout("tcp", host + ":" + port, 1*time.Second)

	if err != nil {
		status = 0
	} else {
		status = 1
		conn.Close()
	}

	results.mutex.Lock()
	defer results.mutex.Unlock()
	fmt.Println(port)
	results.scans = append(results.scans, Scan{status: status, host: host, port:port})
}

func scanAll(host string, ports []string) Results {
	var results Results
	var wg sync.WaitGroup

	wg.Add(len(ports))
	for _, port := range ports {
		go func(p string) {
			defer wg.Done()
			scanPort(host, p, &results)
		}(port)
	}
	wg.Wait()
	fmt.Println(results.scans)
	return results
}

func main() {
	args := os.Args[1:]
	host := args[0]
	ports := parsePorts(args[1])

	// scan host with given ports
	scanAll(host, ports)
}
