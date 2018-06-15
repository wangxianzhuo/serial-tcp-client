package client

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/wangxianzhuo/calc-tool/modbus"
)

// Client modbus client(via tcp)
type Client struct {
	conn     net.Conn
	deadline time.Duration
}

// Connect connect to the tcp server of modbus rtu master
func (c *Client) Connect(serverAddr string, deadline time.Duration) error {
	var err error
	c.conn, err = net.Dial("tcp", serverAddr)
	if err != nil {
		return fmt.Errorf("connect to %v error: %v", serverAddr, err)
	}
	log.Printf("connect to server %v", c.conn.RemoteAddr())
	c.deadline = deadline

	c.conn.SetDeadline(time.Now().Add(deadline))

	return nil
}

// Request execute request
func (c *Client) Request(instruction []byte) ([]byte, error) {
	// write
	_, err := c.conn.Write(instruction)
	if err != nil {
		return nil, fmt.Errorf("send %X to server %v error: %v", instruction, c.conn.RemoteAddr(), err)
	}
	log.Printf("send %X to server %v", instruction, c.conn.RemoteAddr())

	// read
	buf := make([]byte, 100)
	n, err := c.conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("read from server %v error: %v", c.conn.RemoteAddr(), err)
	}

	// prolong dead line
	c.conn.SetDeadline(time.Now().Add(c.deadline))

	// response check
	if err := modbus.ResponseCheck(instruction, buf[:n]); err != nil {
		return nil, err
	}
	log.Printf("read %X from server %v", buf[:n], c.conn.RemoteAddr())

	return buf[:n], nil
}

// Close ...
func (c *Client) Close() {
	c.conn.Close()
}
