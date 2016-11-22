package main

import (
	"encoding/gob"
	"net"
)

type Client struct {
	unixaddr string
}

func (client *Client) send(r *Request) (*Response, error) {
	// Connect
	c, err := net.DialUnix("unix", nil, &net.UnixAddr{client.unixaddr, "unix"})
	if err != nil {
		return nil, err
	}

	// Send request
	if err = gob.NewEncoder(c).Encode(r); err != nil {
		return nil, err
	}
	if err = c.CloseWrite(); err != nil {
		return nil, err
	}

	// Get response
	var resp Response
	err = gob.NewDecoder(c).Decode(&resp)
	if err != nil {
		return nil, err
	}
	if err = c.CloseRead(); err != nil {
		return nil, err
	}

	return &resp, nil
}
