package main

import (
	d "github.com/koalazub/rocket-crash/database"
	_ "github.com/libsql/libsql-client-go/libsql"
)

func main() {
	d.StartDatabase()
}
