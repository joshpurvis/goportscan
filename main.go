package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Scan struct {
	host   string
	port   string
	status int
}

type Result struct {
	mutex sync.Mutex
	scans []Scan
}

func parsePorts(portsRaw string) ([]string, error) {

	var ports []string

	if strings.Contains(portsRaw, ":") {
		// range of ports
		portSplit := strings.Split(portsRaw, ":")

		// parse first port
		firstPort, err := strconv.Atoi(portSplit[0])
		if err != nil {
			return nil, err
		}

		// parse second port
		secondPort, err := strconv.Atoi(portSplit[1])
		if err != nil {
			return nil, err
		}

		for i := firstPort; i <= secondPort; i++ {
			ports = append(ports, strconv.Itoa(i))
		}

	} else {
		// single port
		ports = append(ports, portsRaw)
	}

	return ports, nil
}

func scanPort(host string, port string, result *Result) {
	var status int

	conn, err := net.DialTimeout("tcp", host+":"+port, 1*time.Second)

	if err != nil {
		status = 0
	} else {
		status = 1
		conn.Close()
	}

	result.mutex.Lock()
	defer result.mutex.Unlock()
	result.scans = append(result.scans, Scan{status: status, host: host, port: port})
}

func scanAll(host string, ports []string) Result {
	var result Result
	var wg sync.WaitGroup

	wg.Add(len(ports))
	for _, port := range ports {
		go func(p string) {
			defer wg.Done()
			scanPort(host, p, &result)
		}(port)
	}
	wg.Wait()

	return result
}

func main() {
	var host string

	// parse command line arguments
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Usage: goportscan HOST PORT")
		os.Exit(0)
	}

	host = args[0]
	ports, err := parsePorts(args[1])

	if err != nil {
		fmt.Println("Failed to parse ports")
		os.Exit(1)
	} else {
		// scan host with given ports
		result := scanAll(host, ports)
		fmt.Println(result.scans)
	}

}
