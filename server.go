package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	unixpath string
	paths    Shortcuts
	file     *os.File
}

func NewServer(unixpath string, savefile string) (*Server, error) {
	// Open gob file to load paths
	file, err := os.OpenFile(savefile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	s := &Server{unixpath, NewShortcuts(), file}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if stat.Size() > 0 {
		if err = s.Load(); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s Server) Load() error {
	return gob.NewDecoder(s.file).Decode(&s.paths)
}

func (s Server) Save() error {
	s.file.Seek(0, 0)
	s.file.Truncate(0)

	return gob.NewEncoder(s.file).Encode(s.paths)
}

func (s *Server) Close() {
	if err := s.Save(); err != nil {
		fmt.Println(err)
	}
	s.file.Close()
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

	return strings.Join(matched, "\n"), nil
}

func (s *Server) request(req string) (string, error) {
	// First, return shortcut if exists
	shortcut, err := s.paths.Get(req)
	if err != nil {
		return "", err
	} else if shortcut != "" {
		return shortcut, nil
	}

	// Here, shortcut not exists. Add it to shortcuts.
	// Add func returns absolute req path
	return s.paths.Add(req)
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

	fmt.Println("Request:", req)

	if req.Type == CompletionRequest {
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

	ticker := time.NewTicker(time.Hour)

	// Wait conn or signal
	for run {
		select {
		case <-ticker.C:
			fmt.Println("Save data")
			if err := s.Save(); err != nil {
				fmt.Println(err)
			}
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
