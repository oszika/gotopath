package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

type Server struct {
	unixpath string

	// TODO: Manage same shortcut for several paths
	// map[string][]string
	paths map[string]string
}

func NewServer(unixpath string) *Server {
	return &Server{unixpath, make(map[string]string)}
}

func (s *Server) complete(req string) (string, error) {
	matched := []string{}

	for key, _ := range s.paths {
		ok, err := regexp.MatchString(req, key)
		if err != nil {
			return "", err
		}

		if ok {
			matched = append(matched, key)
		}
	}

	return strings.Join(matched, ","), nil
}

func (s *Server) request(req string) (string, error) {
	fmt.Println("Request:", req)

	// First, return value in paths map
	if resp, ok := s.paths[req]; ok {
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
	s.paths[filepath.Base(req)] = resp
	fmt.Println("Paths:", s.paths)
	fmt.Println("Response:", resp)

	return resp, nil
}

func (s *Server) handleConn(c *net.UnixConn) error {
	defer c.Close()

	// Get request
	var req Request
	err := gob.NewDecoder(c).Decode(&req)
	if err != nil {
		return err
	}
	if err = c.CloseRead(); err != nil {
		return err
	}

	// Treat
	errPath := ""
	var resp string

	if req.Type == Completion {
		resp, err = s.complete(string(req.Req))
		if err != nil {
			errPath = err.Error()
		}
	} else {
		resp, err = s.request(string(req.Req))
		if err != nil {
			errPath = err.Error()
		}
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

func (s *Server) listen() error {
	run := true

	// Signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, os.Kill)
	signal.Notify(signals, syscall.SIGTERM)

	// Conns
	conns := make(chan *net.UnixConn, 100)

	// Listen
	l, err := net.ListenUnix("unix", &net.UnixAddr{s.unixpath, "unix"})
	if err != nil {
		return err
	}

	defer l.Close()
	defer os.Remove(s.unixpath)

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

	fmt.Println("Starting server")

	// Wait conn or signal
	for run {
		select {
		case c := <-conns:
			fmt.Println("Got new conn")
			err := s.handleConn(c)
			if err != nil {
				fmt.Println(err)
			}
		case s := <-signals:
			fmt.Println("Got signal:", s)
			run = false
		}
	}

	fmt.Println("Ending server")

	return nil
}
