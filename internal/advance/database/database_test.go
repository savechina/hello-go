package database

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"gorm.io/gorm"
)

func TestDatabaseFlows(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "crud lifecycle",
			run: func(t *testing.T) {
				db, err := openInMemoryDB()
				if err != nil {
					t.Fatalf("open db: %v", err)
				}
				defer closeDB(db)

				record := learner{Name: "Ada", Email: "ada@example.com"}
				if err := db.Create(&record).Error; err != nil {
					t.Fatalf("create learner: %v", err)
				}

				var loaded learner
				if err := db.First(&loaded, record.ID).Error; err != nil {
					t.Fatalf("read learner: %v", err)
				}
				if loaded.Name != "Ada" {
					t.Fatalf("expected learner name Ada, got %q", loaded.Name)
				}

				if err := db.Model(&loaded).Update("email", "ada@learning.dev").Error; err != nil {
					t.Fatalf("update learner: %v", err)
				}

				var updated learner
				if err := db.First(&updated, record.ID).Error; err != nil {
					t.Fatalf("reload learner: %v", err)
				}
				if updated.Email != "ada@learning.dev" {
					t.Fatalf("expected updated email, got %q", updated.Email)
				}

				if err := db.Delete(&updated).Error; err != nil {
					t.Fatalf("delete learner: %v", err)
				}

				var count int64
				db.Model(&learner{}).Where("id = ?", record.ID).Count(&count)
				if count != 0 {
					t.Fatalf("expected deleted learner count 0, got %d", count)
				}
			},
		},
		{
			name: "preload relationship",
			run: func(t *testing.T) {
				db, err := openInMemoryDB()
				if err != nil {
					t.Fatalf("open db: %v", err)
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
					t.Fatalf("create relationship data: %v", err)
				}

				var loaded learner
				if err := db.Preload("Orders").First(&loaded, record.ID).Error; err != nil {
					t.Fatalf("preload orders: %v", err)
				}
				if len(loaded.Orders) != 2 {
					t.Fatalf("expected 2 orders, got %d", len(loaded.Orders))
				}
			},
		},
		{
			name: "transaction rollback",
			run: func(t *testing.T) {
				db, err := openInMemoryDB()
				if err != nil {
					t.Fatalf("open db: %v", err)
				}
				defer closeDB(db)

				_ = db.Transaction(func(tx *gorm.DB) error {
					if err := tx.Create(&learner{Name: "Rollback", Email: "rollback@example.com"}).Error; err != nil {
						return err
					}

					return errors.New("force rollback")
				})

				var count int64
				db.Model(&learner{}).Where("email = ?", "rollback@example.com").Count(&count)
				if count != 0 {
					t.Fatalf("expected rollback count 0, got %d", count)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func TestRunOutput(t *testing.T) {
	output := captureOutput(t, Run)

	tests := []struct {
		name string
		want string
	}{
		{name: "crud example", want: "示例1 GORM CRUD"},
		{name: "relationship example", want: "示例2 一对多关系"},
		{name: "transaction example", want: "示例3 事务"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(output, tt.want) {
				t.Fatalf("expected output to contain %q, got %q", tt.want, output)
			}
		})
	}
}

func captureOutput(t *testing.T, runner func()) string {
	t.Helper()

	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("create pipe: %v", err)
	}

	os.Stdout = w
	runner()
	_ = w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("read output: %v", err)
	}

	return buf.String()
}
