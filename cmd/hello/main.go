package main

import (
	"fmt"
	"hello/internal/domain"
	"hello/internal/first"
	"hello/internal/repository/boltdb"
	"hello/internal/repository/sqlite"
	"hello/internal/version"
	// "github.com/spf13/cobra"
)

func main() {
	println(version.VERSION)

	first.First()

	// new object
	person := new(domain.Person)
	person.ID = 1
	person.Name = "jack"
	person.Age = 32
	person.Sex = 1

	name := person.GetName()

	println("person:", name)

	//anther new object

	var person2 = &domain.Person{
		ID:   2,
		Name: "Wee",
		Age:  21,
		Sex:  2,
	}

	fmt.Println("person 2 :{}", person2)

	println("repository sqlite3 sample...")

	boltdb.Bolt_demo()

	println("repository sqlite3 sample...")
	sqlite.Sqlite3_demo()

	sqlite.Sqlite3_orm_demo()

}
