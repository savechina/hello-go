package main

import (
	"hello/internal/boltdb"
	"hello/internal/first"
	"hello/internal/version"
)

func main() {
	println(version.VERSION)

	first.First()

	person := new(first.Person)
	person.Name = "jack"
	person.Age = 32
	person.Sex = 1

	name := person.GetName()

	println(name)

	boltdb.Bolt_test()

}
