package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

type client struct {
	connection net.Conn
	inScanner  *bufio.Scanner
	outScanner *bufio.Scanner
	in         io.ReadCloser
	out        io.Writer
	address    string
	timeout    time.Duration
}

func (c client) Close() error {
	if c.connection != nil {
		err := c.connection.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c client) Send() error {
	success := c.inScanner.Scan()
	if !success {
		err := c.Close()
		return err
	}

	message := c.inScanner.Text() + "\n"
	_, err := c.connection.Write([]byte(message))
	return err
}

func (c client) Receive() error {
	success := c.outScanner.Scan()
	if !success {
		err := c.Close()
		return err
	}

	_, err := c.out.Write([]byte(c.outScanner.Text() + "\n"))
	return err
}

func (c *client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		fmt.Println("No connection!")
		return err
	}
	c.connection = conn

	inScanner := bufio.NewScanner(c.in)
	outScanner := bufio.NewScanner(c.connection)

	inScanner.Split(bufio.ScanLines)
	outScanner.Split(bufio.ScanLines)

	c.inScanner = inScanner
	c.outScanner = outScanner

	return nil
}

func createTelnetClient(in io.ReadCloser, out io.Writer, address string, timeout time.Duration) TelnetClient {
	return &client{in: in, out:out, address: address, timeout: timeout}
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return createTelnetClient(in, out, address, timeout)
}
