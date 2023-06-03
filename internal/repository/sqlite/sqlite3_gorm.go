package sqlite

import (
	"fmt"
	"hello/internal/domain"

	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Sqlite3_orm_demo() {
	// Open a connection to a SQLite3 database.
	db, err := openDB()

	if err != nil {
		log.Fatal(err)
	}

	// Create a new user.
	user := domain.User{
		Name: "John Doe",
		Age:  30,
	}

	// Create a table for users.
	db.AutoMigrate(&user)

	db.Create(&user)

	// Get all users.
	users := []domain.User{}
	db.Find(&users)

	// Print all users.
	for _, user := range users {
		fmt.Println(user)
	}

	//gorm read
	var userRepo = &UserRepo{&SqliteRepo{db: db}}

	var u, e = userRepo.FindUsers(nil)

	if e != nil {
		panic(e)
	}

	fmt.Println("userRepo FindUsers:", u)
}

func (userRepo *UserRepo) FindUsers(map[string]interface{}) (*domain.User, error) {
	var user domain.User

	userRepo.db.First(&user)

	return &user, nil
}

// init db conection
func openDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("data/foo.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db, err
}
