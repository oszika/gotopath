package main

// Request type
type Type int

const (
	Path Type = iota
	Completion
)

type Request struct {
	Type Type
	Req  string
}
