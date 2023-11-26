package main

import (
	d "github.com/koalazub/rocket-crash/database"
	s "github.com/koalazub/rocket-crash/server"
	_ "github.com/libsql/libsql-client-go/libsql"
)

func main() {
	db := d.Start()
	s.InitServer(db)

}
