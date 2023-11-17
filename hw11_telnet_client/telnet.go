package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Telneclilient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	address    string
	timeout    time.Duration
	connection net.Conn
	in         io.ReadCloser
	out        io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) Telneclilient {
	return &Client{address: address, timeout: timeout, in: in, out: out}
}

func (cli *Client) Connect() (err error) {
	cli.connection, err = net.DialTimeout("tcp", cli.address, cli.timeout)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", cli.address)
	return nil
}

func (cli *Client) Close() error {
	err := cli.connection.Close()
	if err != nil {
		return err
	}
	return nil
}

func (cli *Client) Receive() error {
	_, err := io.Copy(cli.out, cli.connection)
	fmt.Fprintf(os.Stderr, "...Connection closed by peer\n")
	return err
}

func (cli *Client) Send() error {
	_, err := io.Copy(cli.connection, cli.in)
	fmt.Fprintf(os.Stderr, "...EOF\n")
	return err
}
