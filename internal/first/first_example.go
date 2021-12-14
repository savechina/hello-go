package first

import (
	"fmt"
	"hello/internal/version"
)

func First() {
	sayHello()
}

func sayHello() {
	fmt.Println("hello world!")

	fmt.Println(version.VERSION)
}

/**
Person 类实例使用
*/
type Person struct {
	Name string
	Age  int
	Sex  int
}

func (person *Person) GetName() string {
	return person.Name
}
