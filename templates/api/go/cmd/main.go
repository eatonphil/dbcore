package main

import "fmt"
import "{{ api.extra.repo }}/{{ out_dir }}/pkg/server"

func main() {
	cfg, err := server.NewConfig("../api.yml")
	if err != nil {
		panic(err)
	}

	s, err := server.New(cfg)
	if err != nil {
		panic(err)
	}

	s.Start()
}
