package main

import (
	"net"
)

type Client struct {
	unixaddr string
}

func (client *Client) send(r Request) (*Response, error) {
	// Connect
	c, err := net.DialUnix("unix", nil, &net.UnixAddr{client.unixaddr, "unix"})
	if err != nil {
		return nil, err
	}

	conn := NewConn(c)

	// Send request
	if err = conn.Encode(&r); err != nil {
		return nil, err
	}

	// Get response
	var resp Response
	if err = conn.Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
