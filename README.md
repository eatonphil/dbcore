# DBCore

DBCore is a code generator build around database schemas. Included
with dbcore are templates for generating a Go REST API and a
React/TypeScript browser frontend.

## Example

To build the todo app:

```bash
$ git clone git@github.com:eatonphil/dbcore
$ cd dbcore
$ dotnet run ./examples/todo
$ cd ./examples/todo/go
$ go build ./cmd/main.go
$ ./main
INFO[0000] Starting server at :9090                      pkg=server struct=Server
... in a new window ...
$ curl -X POST -d '{"username": "phil", "password": "phil", "name": "Phil"}' localhost:9090/users/new
{"id":1,"username":"phil","password":"phil","name":"Phil"}
$ curl 'localhost:9090/users?limit=25&offset=0&sortColumn=id&sortOrder=desc' | jq
{
  "total": 1,
  "data": [
    {
      "id": 1,
      "username": "phil",
      "password": "phil",
      "name": "Phil"
    },
  ]
}
```

## Dependencies

* Go
* PostgreSQL
* .NET Core

## Features

### Core

* Read from a PostgreSQL database:
  * Tables (and their columns, primary key, and foreign keys)
* Copy `template/api/<template>` directory, filling out [Liquid-style templates](https://github.com/lunet-io/scriban/blob/master/doc/language.md)

### Go API Features

* YAML-based configuration
* Clean shutdown
* Endpoints and models for:
  * Get, insert, update, delete
  * Get many with filtering, sorting, and pagination 


## Restrictions

There are a bunch of restrictions! Here are a few known ones. You will discover more.

* Only PostgreSQL
* Only tables supported (i.e. no views)
* Only tables within the `public` schema supported
* Only single-column foreign keys supported
* Only Go API templates provided

## Adding a new set of templates

`scripts/post-generate.sh` is a required script. It can be empty.
