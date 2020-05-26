test-todo:
	dotnet run ./examples/todo
	cd ./examples/todo/go && make && ./main
