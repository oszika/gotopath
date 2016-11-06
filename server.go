package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
)

func request(req string) string {
	fmt.Println("Request:", req)
	return req
}

func handleConn(c *net.UnixConn) error {
	defer c.Close()

	// Get request
	req, err := ioutil.ReadAll(c)
	if err != nil {
		return err
	}
	if err = c.CloseRead(); err != nil {
		return err
	}

	// Treat
	resp := request(string(req))

	// Send response
	_, err = fmt.Fprintf(c, resp)
	if err != nil {
		return err
	}
	if err = c.CloseWrite(); err != nil {
		return err
	}

	return nil
}

func listen(unixpath string) error {
	run := true

	// Signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, os.Kill)

	// Conns
	conns := make(chan *net.UnixConn, 100)

	// Listen
	l, err := net.ListenUnix("unix", &net.UnixAddr{unixpath, "unix"})
	if err != nil {
		return err
	}

	defer l.Close()
	defer os.Remove(unixpath)

	// Listen connections and send them to conns chan
	go func() {
		for run {
			c, err := l.AcceptUnix()
			if err != nil {
				fmt.Println(err)
				continue
			}

			conns <- c
		}
	}()

	// Wait conn or signal
	for run {
		select {
		case c := <-conns:
			fmt.Println("Got new conn")
			err := handleConn(c)
			if err != nil {
				fmt.Println(err)
			}
		case s := <-signals:
			fmt.Println("Got signal:", s)
			run = false
		}
	}

	return nil
}
