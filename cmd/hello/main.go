package main

import (
	"hello/internal/first"
)

func main() {

	first.First()

	// println(first.VER)

	person := new(first.Person)
	person.Name = "jack"
	person.Age = 32
	person.Sex = 1

	name := person.GetName()

	println(name)

}
