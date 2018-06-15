package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/wangxianzhuo/calc-tool/modbus"

	"github.com/wangxianzhuo/serial-tcp-client/client"
)

var instruction = flag.String("ins", "", "instruction for modbus-rtu")
var instructionWithoutCheckCode = flag.String("ins-withou-check-code", "", "instruction for modbus-rtu without check code")
var tcpAddr = flag.String("tcp-addr", ":9000", "tcp server address")

func main() {
	checkParams()

	// get instruction
	ins, err := getInstructions(*instruction, *instructionWithoutCheckCode)
	if err != nil {
		log.Fatal(err)
	}

	client := client.Client{}
	err = client.Connect(*tcpAddr, time.Second*time.Duration(10))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	_, err = client.Request(ins)
	if err != nil {
		log.Fatal(err)
	}
}

func checkParams() {
	flag.Parse()

	if *instruction == "" && *instructionWithoutCheckCode == "" {
		fmt.Println("Error: --ins or --ins-withou-check-code at least need one")
		fmt.Println("\t--ins", flag.Lookup("ins").Usage)
		fmt.Println("\t--ins-withou-check-code", flag.Lookup("ins-withou-check-code").Usage)
		os.Exit(2)
	}
}

func getInstructions(ins, insWithoutCheckCode string) (instruction []byte, err error) {
	if insWithoutCheckCode != "" {
		instruction, err = modbus.InstructionWithCheckCode(insWithoutCheckCode)
	} else if ins != "" {
		instruction, err = modbus.InstructionParse(ins)
	} else {
		return nil, fmt.Errorf("need insturction")
	}

	if !modbus.CRC16Check(instruction) {
		return nil, fmt.Errorf("instruction %X check failed", instruction)
	}
	return instruction, nil
}
