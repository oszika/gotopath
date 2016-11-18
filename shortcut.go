package main

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
