package main

type Path struct {
	Name  string
	Count int
}

func NewPath(name string) *Path {
	return &Path{name, 1}
}

type Shortcut struct {
	Name string

	Path *Path
	// TODO: manage several paths
	// paths []*Path
}

func NewShortcut(name string, path string) *Shortcut {
	return &Shortcut{name, NewPath(path)}
}
