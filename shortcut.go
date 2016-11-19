package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

type Shortcut struct {
	Main  string         // Main path
	Paths map[string]int // All paths managed by shortcut
}

func NewShortcut(path string) *Shortcut {
	return &Shortcut{path, map[string]int{path: 1}}
}

func (s *Shortcut) Update(path string) {
	if _, ok := s.Paths[path]; ok {
		s.Paths[path]++
		if s.Paths[path] > s.Paths[s.Main] {
			s.Main = path
		}
	} else {
		s.Paths[path] = 1
	}
}

type Shortcuts map[string]*Shortcut

func NewShortcuts() Shortcuts {
	return make(map[string]*Shortcut)
}

func (s Shortcuts) Get(req string) string {
	if shortcut, ok := s[req]; ok {
		fmt.Println("Shortcut found:", shortcut)
		return shortcut.Main
	}

	return ""
}

func (s Shortcuts) Add(req string) (string, error) {
	// Check path
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

	// Add shortcut for each subfile
	for resp != "/" {
		d, b := path.Dir(resp), path.Base(resp)
		if _, ok := s[b]; !ok {
			s[b] = NewShortcut(resp)
			fmt.Println("New shortcut created:", b, "->", s[b])
		} else {
			s[b].Update(resp)
			fmt.Println("Shortcut updated:", b, "->", s[b])
		}
		resp = d
	}

	return resp, nil
}
