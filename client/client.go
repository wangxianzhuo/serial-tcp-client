package client

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/wangxianzhuo/calc-tool/modbus"
)

// Client modbus client(via tcp)
type Client struct{}

// Run execute request
func (c *Client) Run(serverAddr string, instruction []byte) error {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return fmt.Errorf("connect to %v error: %v", serverAddr, err)
	}
	defer conn.Close()
	log.Printf("connect to server %v", conn.RemoteAddr())
	conn.SetDeadline(time.Now().Add(time.Second * time.Duration(10)))

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
	if err := modbus.ResponseCheck(instruction, buf[:n]); err != nil {
		return err
	}
	log.Printf("read %X from server %v", buf[:n], conn.RemoteAddr())

	return nil
}
