package main

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path/filepath"
)

// TODO: Manage same shortcut for several paths
// map[string][]string
var paths map[string]string = make(map[string]string)

func request(req string) (string, error) {
	fmt.Println("Request:", req)

	// First, return value in paths map
	if resp, ok := paths[req]; ok {
		fmt.Println("Response:", resp)
		return resp, nil
	}

	// Check path and add to paths maps
	info, err := os.Stat(req)
	if err != nil {
		return "", err
	}

	// Path must be valid dir
	if !info.IsDir() {
		return "", os.ErrNotExist
	}

	// Add to paths map
	resp, err := filepath.Abs(req)
	if err != nil {
		return "", err

	}
	paths[filepath.Base(req)] = resp
	fmt.Println("Paths:", paths)
	fmt.Println("Response:", resp)

	return resp, nil
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
	errPath := ""
	resp, err := request(string(req))
	if err != nil {
		errPath = err.Error()
	}

	// Send response
	err = gob.NewEncoder(c).Encode(Response{resp, errPath})
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
