package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/akamensky/argparse"
)

type address struct {
	ip   *string
	port int
}

type storeDataField struct {
	port int
}

type storeData struct {
	storeDataField []int
}

var (
	ip *string
	i  *string
)

func init() {

	parser := argparse.NewParser("slow-scanner", "A slow port scanner")
	i = parser.String("i", "ip", &argparse.Options{Required: true, Help: "IP address to scan"})

	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Println(parser.Usage(err))
		return
	}
}

func scan(ch chan<- storeData, ip *string, port int) {

	address := fmt.Sprintf("%s:%d", *ip, port)
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)

	if err != nil {
		return
	}

	ch <- storeData{
		storeDataField: []int{port},
	}

	conn.Close()
}

func main() {

	ip = i
	fmt.Println("IP:", *ip)

	dataChannel := make(chan storeData)

	// Execute port scanner
	for i := 0; i < 1024; i++ {
		go scan(dataChannel, ip, i)
	}

	// Store response from channel to struct
	for i := 0; i < 1024; i++ {

		result := <-dataChannel

		for _, port := range result.storeDataField {
			fmt.Println("\033[1;32m[+]\033[0m Port open:", port)
			if port == 80 {
				close(dataChannel)
				os.Exit(0)
			}

		}
	}
}
