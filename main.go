package main

import (
	"os"
	"fmt"
	"net"
	"time"
	"strings"
	"strconv"
)

type result struct {
	host string
	port string
	status int
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

func scanPort(host string, port string) result {
	var status int

	conn, err := net.DialTimeout("tcp", host + ":" + port, 1*time.Second)

	if err != nil {
		status = 0
	} else {
		status = 1
		conn.Close()
	}

	return result{status: status, host: host, port:port}
}

func scanAll(host string, ports []string) []result {
	var results []result

	for _, port := range ports {
		results = append(results, scanPort(host, port))
	}

	fmt.Println(results)
	return results
}

func main() {
	args := os.Args[1:]
	host := args[0]
	ports := parsePorts(args[1])

	// scan host with given ports
	scanAll(host, ports)
}
