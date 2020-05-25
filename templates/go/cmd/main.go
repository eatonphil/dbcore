package main

import "{{api.repo}}/pkg/server"

func main() {
	s := server.New(server.NewConfig("/etc/{{api.project}}.yml"))
	s.Start()
}
