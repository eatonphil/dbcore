# Genapp

Genapp is a code generator build around database schemas. Applications
are configured with a `genapp.yml` file in the project directory and
existing database tables. Genapp will generate an API for you.

## Dependencies

* Go
* PostgreSQL
* .NET Core

## Example

To build the todo app:

```bash
$ git clone git@github.com:eatonphil/genapp
$ cd genapp
$ dotnet run ./examples/todo
$ cd ./examples/todo/go
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

## Restrictions

There are a bunch of restrictions! Here are a few known ones. You will discover more.

* Only PostgreSQL
* Only tables supported (i.e. no views)
* Only tables within the `public` schema supported
* Only single-column foreign keys supported
* Only Go API templates provided

## Adding a new set of templates

`scripts/post-generate.sh` is a required script. It can be empty.
