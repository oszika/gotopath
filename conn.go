package main

import (
	"encoding/gob"
	"net"
)

type Conn struct {
	conn    *net.UnixConn
	encoder *gob.Encoder
	decoder *gob.Decoder
}

func NewConn(c *net.UnixConn) *Conn {
	return &Conn{c, gob.NewEncoder(c), gob.NewDecoder(c)}
}

func (c *Conn) Encode(data interface{}) error {
	if err := c.encoder.Encode(data); err != nil {
		return err
	}
	if err := c.conn.CloseWrite(); err != nil {
		return err
	}

	return nil
}

func (c *Conn) Decode(data interface{}) error {
	if err := c.decoder.Decode(data); err != nil {
		return err
	}
	if err := c.conn.CloseRead(); err != nil {
		return err
	}

	return nil
}
