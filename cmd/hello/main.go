package main

import (
	"hello/internal/first"
	"hello/internal/repository/boltdb"
	"hello/internal/repository/sqlite"
	"hello/internal/version"
	// "github.com/spf13/cobra"
)

func main() {

	println(version.VERSION)

	first.First()

	println("\nrepository sqlite3 sample...")

	boltdb.Bolt_demo()

	println("\nrepository sqlite3 sample...")
	sqlite.Sqlite3_demo()

	sqlite.Sqlite3_orm_demo()

}
