package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Path struct {
	Name  string
	Count int
}

func NewPath(name string) *Path {
	return &Path{name, 1}
}

type Shortcut struct {
	Path *Path
	// TODO: manage several paths
	// paths []*Path
}

func NewShortcut(path string) *Shortcut {
	return &Shortcut{NewPath(path)}
}

type Shortcuts map[string]*Shortcut

func NewShortcuts() Shortcuts {
	return make(map[string]*Shortcut)
}

func (s Shortcuts) Get(req string) string {
	if shortcut, ok := s[req]; ok {
		return shortcut.Path.Name
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
	s[filepath.Base(req)] = NewShortcut(resp)
	fmt.Println("Paths:", s)
	fmt.Println("Response:", resp)

	return resp, nil

}
