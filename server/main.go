package main

import (
	"flag"
	"fmt"
	"os"

	d "github.com/koalazub/rocket-crash/database"
	s "github.com/koalazub/rocket-crash/server"
	_ "github.com/libsql/libsql-client-go/libsql"
)

func main() {
	toLog := flag.Bool("l", false, "Enable logging mode. Made for http requests and such")
	toHelp := flag.Bool("h", false, "Display help")
	flag.Parse()

	if *toHelp {
		brrHelp()
	}

	s.ToLog = toLog // Logging enabled
	d.ToLog = toLog

	db := d.Start()
	s.InitServer(toLog, db)

}

func brrHelp() {
	fmt.Println("\t-h", "help", "\tprovide help. I know - seems fucking obvious by now hey?")
	fmt.Println("\t-l", "log", "\t\tenable logging. This is to track http requests. Especially because we use http3 and QUIC")
	os.Exit(0)
}
