package main

import (
	"flag"

	d "github.com/koalazub/rocket-crash/database"
	s "github.com/koalazub/rocket-crash/server"
	_ "github.com/libsql/libsql-client-go/libsql"
)

func main() {
	toLog := flag.Bool("l", false, "Enable logging mode. Made for http requests and such")
	flag.Parse()

	s.ToLog = toLog // Logging enabled
	d.ToLog = toLog

	db := d.Start()
	s.InitServer(toLog, db)

}
