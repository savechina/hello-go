package domain

/*
*
Person 类实例使用
*/
type Person struct {
	ID   uint
	Name string
	Age  int
	Sex  int
}

func (person *Person) GetName() string {
	return person.Name
}
