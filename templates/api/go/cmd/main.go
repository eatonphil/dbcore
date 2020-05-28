package main

import "fmt"
import "{{api.extra.repo}}/go/pkg/server"

func main() {
	cfg, err := server.NewConfig("../dbcore.yml")
	if err != nil {
		panic(err)
	}

	s, err := server.New(cfg)
	if err != nil {
		panic(err)
	}

	s.Start()
}
