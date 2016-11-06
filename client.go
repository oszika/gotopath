package main

import (
	"fmt"
	"io/ioutil"
	"net"
)

func clientReq(unixaddr string, path string) error {
	c, err := net.DialUnix("unix", nil, &net.UnixAddr{unixaddr, "unix"})
	if err != nil {
		return err
	}

	// Make request
	if _, err = fmt.Fprintf(c, path); err != nil {
		return err
	}
	if err = c.CloseWrite(); err != nil {
		return err
	}

	// Get response
	req, err := ioutil.ReadAll(c)
	if err != nil {
		return err
	}
	if err = c.CloseRead(); err != nil {
		return err
	}

	fmt.Println("Response:", string(req))

	return nil
}
