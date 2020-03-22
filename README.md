# Test SQLite CRUD operations with Goalng

Use this CLI example to learn how to use Golang with SQLite. Modify and reuse some parts of this to make a web app etc.

## How to prepare the project before you test it

First, help the models packages ready to work with main package.

```console
$mv models && go build
```

You can also optionally set up live edit environment with nodemon.

```console
$yarn
```

Then, you can use $./go.sh instead of manually typing $go run main.go

## How to setup SQLite database

You can manually handle SQLite database. [Refer to the documenation for the SQLite CLI.](https://sqlite.org/cli.html)

```console
$touch users.db && sqlite3 users.db
$CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT NOT NULL UNIQUE)
```

You can also uncomment the code similar to this in main.go

```console
db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT NOT NULL UNIQUE)")
```

## How to test it

Use one of commands below to test CRULD(Create, Read, Update, List, Delete) users.

```console
$go run main.go -action=create
$go run main.go -action=get
$go run main.go -action=update
$go run main.go -action=delete
$go run main.go -action=list
```
