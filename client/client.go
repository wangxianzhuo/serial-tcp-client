package client

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/wangxianzhuo/serial-tcp-client/util"
)

type Client struct{}

func (c *Client) Run(serverAddr, ins string, isCrcGen bool) error {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return fmt.Errorf("connect to %v error: %v", serverAddr, err)
	}
	defer conn.Close()
	log.Printf("connect to server %v", conn.RemoteAddr())
	conn.SetDeadline(time.Now().Add(time.Second * time.Duration(10)))

	instruction, err := util.ParseInstruction(ins)
	if err != nil {
		return fmt.Errorf("parse instruction %X error: %v", instruction, err)
	}
	if isCrcGen {
		hi, lo, _ := util.CRC16ModbusCheckCode(instruction)
		instruction = append(instruction, hi)
		instruction = append(instruction, lo)
	}

	_, err = conn.Write(instruction)
	if err != nil {
		return fmt.Errorf("send %X to server %v error: %v", instruction, conn.RemoteAddr(), err)
	}
	log.Printf("send %X to server %v", instruction, conn.RemoteAddr())

	buf := make([]byte, 100)
	n, err := conn.Read(buf)
	if err != nil {
		return fmt.Errorf("read from server %v error: %v", conn.RemoteAddr(), err)
	}
	if !util.CRC16ModbusCheck(buf[:n]) {
		return fmt.Errorf("check sum for %X illegal", buf[:n])
	}
	log.Printf("read %X from server %v", buf[:n], conn.RemoteAddr())

	return nil
}
