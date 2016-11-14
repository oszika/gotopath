package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func requestOrPwd(req string) string {
	if req != "" {
		return req
	}
	return os.Getenv("PWD")
}

func main() {
	serve := flag.Bool("serve", false, "Serve mode")
	completion := flag.Bool("complete", false, "Get suggestions for uncomplete request")
	request := flag.String("request", "", "Where you want to go (path or shortcut)")
	flag.Parse()

	unixaddr := "/tmp/gotopath." + os.Getenv("USER")

	if *serve {
		s, err := NewServer(unixaddr, "/tmp/gotopath.save")
		if err != nil {
			panic(err)
		}
		defer s.Close()
		if err = s.listen(); err != nil {
			panic(err)
		}
	} else {
		var req *Request

		if *completion {
			req = &Request{Completion, *request}
		} else {
			req = &Request{Path, requestOrPwd(*request)}
		}

		resp, err := (&Client{unixaddr}).send(req)
		if err != nil {
			panic(err)
		}
		if resp.Err != "" {
			panic(errors.New(resp.Err))
		}

		// Display result
		fmt.Println(resp.Path)
	}
}
