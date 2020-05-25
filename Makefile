test-todo:
	rm -rf ./examples/todo/go
	dotnet run ./examples/todo
	cd ./examples/todo/go && make
