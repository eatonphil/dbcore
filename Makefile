.PHONY: clean test-todo

build: ./bin/Debug/netcoreapp3.0/dbcore
	dotnet build

clean:
	dotnet clean

test-todo:
	dotnet run ./examples/todo
	cd ./examples/todo/go && make && ./main
