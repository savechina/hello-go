package domain

type User struct {
	ID   uint `gorm:"primary_key"`
	Name string
	Age  int
}
