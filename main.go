package main

import (
	"flag"
	"os"
)

func main() {
	serve := flag.Bool("serve", false, "Serve mode")
	flag.Parse()

	unixaddr := "/tmp/gotopath." + os.Getenv("USER")

	if *serve {
		if err := listen(unixaddr); err != nil {
			panic(err)
		}
	} else {
		var path string

		if len(os.Args) > 1 {
			path = os.Args[1]
		} else {
			path = os.Getenv("PWD")
		}

		if err := clientReq(unixaddr, path); err != nil {
			panic(err)
		}
	}
}
