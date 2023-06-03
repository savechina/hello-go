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
