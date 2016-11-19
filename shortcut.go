package main

import (
	"fmt"
	"os"
	"path"
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
		fmt.Println("Shortcut found:", shortcut)
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

	// Add shortcut for each subfile
	for resp != "/" {
		d, b := path.Dir(resp), path.Base(resp)
		if _, ok := s[b]; !ok {
			s[b] = NewShortcut(resp)
			fmt.Println("New shortcut created:", b, "->", resp)
		} else {
			fmt.Println("Shortcut updated:", b, "->", resp)
		}
		resp = d
	}

	s[filepath.Base(req)] = NewShortcut(resp)
	fmt.Println("Paths:", s)
	fmt.Println("Response:", resp)

	return resp, nil

}
