package database

import (
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"hello/internal/chapters"
)

func init() {
	chapters.Register("advance", "database", Run)
}

// Run prints runnable examples for GORM models, migrations, CRUD, relationships, and transactions.
func Run() {
	examples := []string{
		crudExample(),
		relationshipExample(),
		transactionExample(),
	}

	for _, example := range examples {
		fmt.Println(example)
	}
}

type learner struct {
	ID     uint
	Name   string
	Email  string
	Orders []purchase `gorm:"constraint:OnDelete:CASCADE;"`
}

type purchase struct {
	ID         uint
	LearnerID  uint
	Item       string
	TotalCents int
}

func crudExample() string {
	db, err := openInMemoryDB()
	if err != nil {
		return fmt.Sprintf("示例1 GORM CRUD: error=%v", err)
	}
	defer closeDB(db)

	record := learner{Name: "Ada", Email: "ada@example.com"}
	if err := db.Create(&record).Error; err != nil {
		return fmt.Sprintf("示例1 GORM CRUD: create error=%v", err)
	}

	var loaded learner
	if err := db.First(&loaded, record.ID).Error; err != nil {
		return fmt.Sprintf("示例1 GORM CRUD: read error=%v", err)
	}

	if err := db.Model(&loaded).Update("email", "ada@learning.dev").Error; err != nil {
		return fmt.Sprintf("示例1 GORM CRUD: update error=%v", err)
	}

	var updated learner
	if err := db.First(&updated, record.ID).Error; err != nil {
		return fmt.Sprintf("示例1 GORM CRUD: verify update error=%v", err)
	}

	if err := db.Delete(&updated).Error; err != nil {
		return fmt.Sprintf("示例1 GORM CRUD: delete error=%v", err)
	}

	var remaining int64
	db.Model(&learner{}).Where("id = ?", record.ID).Count(&remaining)

	return fmt.Sprintf(
		"示例1 GORM CRUD: id=%d name=%s updated_email=%s remaining=%d",
		record.ID,
		loaded.Name,
		updated.Email,
		remaining,
	)
}

func relationshipExample() string {
	db, err := openInMemoryDB()
	if err != nil {
		return fmt.Sprintf("示例2 一对多关系: error=%v", err)
	}
	defer closeDB(db)

	record := learner{
		Name:  "Grace",
		Email: "grace@example.com",
		Orders: []purchase{
			{Item: "Go Web", TotalCents: 3590},
			{Item: "GORM Deep Dive", TotalCents: 4290},
		},
	}

	if err := db.Create(&record).Error; err != nil {
		return fmt.Sprintf("示例2 一对多关系: create error=%v", err)
	}

	var loaded learner
	if err := db.Preload("Orders").First(&loaded, record.ID).Error; err != nil {
		return fmt.Sprintf("示例2 一对多关系: preload error=%v", err)
	}

	totalCents := 0
	for _, order := range loaded.Orders {
		totalCents += order.TotalCents
	}

	return fmt.Sprintf(
		"示例2 一对多关系: learner=%s orders=%d total_cents=%d first_item=%s",
		loaded.Name,
		len(loaded.Orders),
		totalCents,
		loaded.Orders[0].Item,
	)
}

func transactionExample() string {
	db, err := openInMemoryDB()
	if err != nil {
		return fmt.Sprintf("示例3 事务: error=%v", err)
	}
	defer closeDB(db)

	if err := db.Transaction(func(tx *gorm.DB) error {
		user := learner{Name: "Linus", Email: "linus@example.com"}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		return tx.Create(&purchase{LearnerID: user.ID, Item: "Transaction Lab", TotalCents: 4990}).Error
	}); err != nil {
		return fmt.Sprintf("示例3 事务: commit error=%v", err)
	}

	rollbackErr := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&learner{Name: "Rollback", Email: "rollback@example.com"}).Error; err != nil {
			return err
		}

		return errors.New("force rollback")
	})

	var committed int64
	db.Model(&learner{}).Where("email = ?", "linus@example.com").Count(&committed)

	var rolledBack int64
	db.Model(&learner{}).Where("email = ?", "rollback@example.com").Count(&rolledBack)

	return fmt.Sprintf(
		"示例3 事务: commit_count=%d rollback_count=%d rollback_err=%t",
		committed,
		rolledBack,
		rollbackErr != nil,
	)
}

func openInMemoryDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(1)

	if err := db.AutoMigrate(&learner{}, &purchase{}); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	return db, nil
}

func closeDB(db *gorm.DB) {
	if db == nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		return
	}

	_ = sqlDB.Close()
}
