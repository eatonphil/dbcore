package main

import "{{repo}}/pkg/server"

func main() {
	s := server.New(server.NewConfig("/etc/{{project}}.yml"))
	s.Start()
}
