package main

import (
	"fmt"
	"os"
	"path"
)

// pathIsValid checks is path is a valid directory.
func pathIsValid(req string) error {
	if !path.IsAbs(req) {
		return os.ErrInvalid
	}
	info, err := os.Stat(req)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return os.ErrNotExist
	}
	return nil
}

// A shortcuts contains several paths with same base. The most used path
// is considerated as 'Main'.
type Shortcut struct {
	Main  string         // Main path
	Paths map[string]int // All paths managed by shortcut
}

func NewShortcut(path string) *Shortcut {
	return &Shortcut{path, map[string]int{path: 1}}
}

// Update updates count accesses to path if exists or adds it.
// The Main is updated if necessary.
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

// Removes a path from Shortcut. If this path is the Main, main is reaffected
// to the new most used path, or "" if not exists.
func (s *Shortcut) Remove(path string) {
	delete(s.Paths, path)

	if path == s.Main {
		c := 0
		s.Main = ""
		for k, v := range s.Paths {
			if v > c {
				s.Main = k
				c = v
			}
		}
	}
}

// Returns if shortcut contains paths
func (s *Shortcut) IsEmpty() bool {
	return len(s.Paths) == 0
}

type Shortcuts map[string]*Shortcut

func NewShortcuts() Shortcuts {
	return make(map[string]*Shortcut)
}

func (s Shortcuts) Remove(base string, path string) bool {
	// Here, path is absolute
	if shortcut, ok := s[base]; ok {
		fmt.Println("Remove:", base, "->", path)
		shortcut.Remove(path)
		if shortcut.IsEmpty() {
			fmt.Println("Delete:", base)
			delete(s, base)
			return true
		}
	}
	return false
}

func (s Shortcuts) RemoveAllInvalidPaths(req string) {
	// req is absolute
	for req != "/" {
		d, b := path.Dir(req), path.Base(req)
		if err := pathIsValid(req); err != nil {
			s.Remove(b, req)
		}
		req = d
	}

}

func (s Shortcuts) Get(req string) string {
	fmt.Println("Get")
	if shortcut, ok := s[req]; ok {
		// Execute for all paths in shortcut
		for !shortcut.IsEmpty() {
			if err := pathIsValid(shortcut.Main); err == nil {
				// Update shortcut
				s.Update(shortcut.Main)

				fmt.Println("Shortcut updated:", shortcut)

				return shortcut.Main
			} else {
				// Path not valid, remove it and check the next
				// If there is only one paths, shortcut will be destroy
				if del := s.Remove(req, shortcut.Main); del {
					return ""
				}
			}

		}

		return ""
	}

	return ""
}

func (s Shortcuts) Add(req string) (string, error) {
	err := pathIsValid(req)
	if err != nil {
		s.RemoveAllInvalidPaths(req)
		return "", err
	}

	// Add shortcut for each subfile
	s.Update(req)

	return req, nil
}

func (s Shortcuts) Update(req string) {
	for req != "/" {
		d, b := path.Dir(req), path.Base(req)
		if _, ok := s[b]; !ok {
			s[b] = NewShortcut(req)
			fmt.Println("New shortcut created:", b, "->", s[b])
		} else {
			s[b].Update(req)
			fmt.Println("Shortcut updated:", b, "->", s[b])
		}
		req = d
	}
}
