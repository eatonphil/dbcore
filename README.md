# Genapp

To build the todo app:

```bash
dotnet run ./examples/todo
```

## Restrictions

There are a bunch of restrictions! Here are a few known ones. You will discover more.

* Only tables supported (i.e. no views)
* Only tables within the `public` schema supported
* Only single-column foreign keys supported
* Only Go API templates provided

## Adding a new set of templates

`scripts/post-generate.sh` is a required script. It can be empty.
