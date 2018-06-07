package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wangxianzhuo/serial-tcp-client/client"
)

var instruction = flag.String("ins", "", "instruction for modbus-rtu")
var tcpAddr = flag.String("tcp-addr", ":9000", "tcp server address")
var crcGenerate = flag.Bool("crc-gen", false, "can generate crc check code")

func main() {
	flag.Parse()

	checkParams()

	client := client.Client{}
	err := client.Run(*tcpAddr, *instruction, *crcGenerate)
	if err != nil {
		panic(err)
	}
}

func checkParams() {
	if *instruction == "" {
		fmt.Println("Error: ins can't be empty")
		flag.PrintDefaults()
		os.Exit(2)
	}
}
