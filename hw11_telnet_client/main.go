package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "Timeout for connection")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, errors.New("invalid arguments, must specify host and port"))
		os.Exit(1)
	}

	client := NewTelnetClient(
		net.JoinHostPort(args[0], args[1]),
		*timeout,
		os.Stdin,
		os.Stdout,
	)
	err := client.Connect()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer client.Close()

	go func() {
		defer cancel()
		err := client.Receive()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}()

	go func() {
		defer cancel()
		err := client.Send()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}()

	<-ctx.Done()
}
