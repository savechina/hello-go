/*
包
每个 Go 程序都是由包构成的。
*/
package first

/* 按照约定，包名与导入路径的最后一个元素一致。
例如，"math/rand" 包中的源码均以 package rand 语句开始。
*/
import (
	"fmt"
	"hello/internal/domain"
	"hello/internal/version"
)

/*
Hello Example  First
*/
func First() {

	fmt.Println("\nHE01：Hello World \nOutput:")
	sayHello()

	fmt.Println("\nHE02：短变量声明 example \nOutput:")
	sayVariable()

	// new object
	//
	fmt.Println("\nHE03：创建对象 example \nOutput:")
	sayNewObject()

}

/*
HE01：hello world
*/
func sayHello() {

	//打印输出
	fmt.Println("hello world!")

	//打印版本号：
	fmt.Println(version.VERSION)
}

/*
HE02：短变量声明 example
*/
func sayVariable() {
	// 短变量声明
	// 在函数中，简洁赋值语句 := 可在类型明确的地方代替 var 声明。
	var i, j int = 1, 2
	k := 3
	c, python, java := true, false, "no!"

	fmt.Println(i, j, k, c, python, java)
}

/*
HE03：创建对象
*/
func sayNewObject() {
	//使用 new 创建对象
	person := new(domain.Person)
	person.ID = 1
	person.Name = "jack"
	person.Age = 32
	person.Sex = 1

	name := person.GetName()

	println("new person [new(Person)]:", name)

	// anther new object
	//直接创建对象
	var person2 = &domain.Person{
		ID:   2,
		Name: "Wee",
		Age:  21,
		Sex:  2,
	}

	fmt.Println("new person [&Person{}]:", person2)
}
