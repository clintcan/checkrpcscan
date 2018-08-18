package main

/*
	checkrpc.go - this simple go program will check for open rpc geth interfaces
	Program by: Clint Canada
	Usage:
		checkrpc.exe -i <ip address default 127.0.0.1> -p <port, default 8545>
*/

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func ScanPort(ip string, port int) bool {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, 2*time.Second)

	if err != nil {
		return false
	}

	conn.Close()
	return true
}

func CheckRPC(ip string, port int) bool {
	address := "http://" + ip + ":" + strconv.Itoa(port)
	var jsonStr = []byte("{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[],\"id\":1}")
	req, err := http.NewRequest("POST", address, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	testFlag := strings.Contains(string(body), "Geth")
	return testFlag
}

func main() {
	// We shall get the ip address and the port
	ipaddr := flag.String("i", "127.0.0.1", "ip address to check")
	portaddr := flag.Int("p", 8545, "geth port to check if open")

	flag.Parse()

	ip := *ipaddr
	port := *portaddr

	portstr := strconv.Itoa(port)

	fmt.Println("Scanning " + ip + " port " + portstr)

	portstatus := ScanPort(ip, port)

	if portstatus == true {
		fmt.Println(port, "open, checking for misconfigured RPC Port...")
		testFlag := CheckRPC(ip, port)
		if testFlag == true {
			fmt.Println("RPC Port is found.")
		} else {
			fmt.Println("No RPC interface found.")
		}
		os.Exit(0)
	}

	fmt.Println(port, "closed, exiting")
}
