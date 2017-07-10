package main

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type Request interface {
	cb(Shortcuts) (string, error)
}

type RequestCompletion struct {
	To string
}

func (r RequestCompletion) cb(paths Shortcuts) (string, error) {
	matched := []string{}

	for key, shortcut := range paths {
		ok, err := regexp.MatchString(r.To, key)
		if err != nil {
			return "", err
		}

		if ok {
			matched = append(matched, key)
		}

		for path, _ := range shortcut.Paths {
			ok, err := regexp.MatchString(r.To, path)
			if err != nil {
				return "", err
			}

			if ok {
				matched = append(matched, key+":="+path)
			}

		}
	}

	return strings.Join(matched, "\n"), nil
}

type RequestPath struct {
	To   string
	From string
}

func (r RequestPath) cb(paths Shortcuts) (string, error) {
	// For zsh completion, request can have format: "<shortcut>:=<path>"
	// Get real request
	chunks := strings.Split(r.To, ":=")
	if len(chunks) == 2 {
		r.To = chunks[1]
	}

	// Return shortcut if exists
	if shortcut := paths.Get(r.To); shortcut != "" {
		return shortcut, nil
	}

	// Build complete path if it's not absolute
	if !path.IsAbs(r.To) {
		if !path.IsAbs(r.From) {
			return "", os.ErrInvalid
		}

		r.To = r.From + "/" + r.To
	}

	// Clean path
	var err error
	r.To, err = filepath.Abs(r.To)
	if err != nil {
		return "", err
	}

	// Here, shortcut not exists. Add it to shortcuts.
	// Add func returns absolute req path
	return paths.Add(r.To)
}
