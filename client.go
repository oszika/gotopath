package main

import (
	"encoding/gob"
	"fmt"
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
	var resp Response
	err = gob.NewDecoder(c).Decode(&resp)
	if err != nil {
		return err
	}
	if err = c.CloseRead(); err != nil {
		return err
	}

	fmt.Println("Response:", resp)

	return nil
}
