package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
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

func listen() {
	// Listen
	path := "/tmp/gotopath." + os.Getenv("USER")
	l, err := net.ListenUnix("unix", &net.UnixAddr{path, "unix"})
	if err != nil {
		panic(err)
	}

	defer l.Close()
	defer os.Remove(path)

	for {
		c, err := l.AcceptUnix()
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = handleConn(c)
		if err != nil {
			fmt.Println(err)
		}
	}
}
