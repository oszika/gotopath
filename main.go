package main

import (
	"flag"
	"os"
)

func main() {
	serve := flag.Bool("serve", false, "Serve mode")
	completion := flag.Bool("complete", false, "Get suggestions for uncomplete request")
	flag.Parse()

	unixaddr := "/tmp/gotopath." + os.Getenv("USER")

	if *serve {
		if err := listen(unixaddr); err != nil {
			panic(err)
		}
	} else if *completion {
		var req string

		if len(os.Args) > 1 {
			req = os.Args[1]
		}

		if err := completionReq(unixaddr, req); err != nil {
			panic(err)
		}
	} else {
		var path string

		if len(os.Args) > 1 {
			path = os.Args[1]
		} else {
			path = os.Getenv("PWD")
		}

		if err := pathReq(unixaddr, path); err != nil {
			panic(err)
		}
	}
}
