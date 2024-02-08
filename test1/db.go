package test1

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func RunExport() {
	db, err := gorm.Open(sqlite.Open("clipboard.db"), &gorm.Config{})
	if err != nil {
		fmt.Printf("error connecting to database: %v\n", err)
		return
	}

	err = db.AutoMigrate(&ClipboardEntry{})
	if err != nil {
		fmt.Printf("error migrating database: %v\n", err)
		return
	}

	if err := exportToJSON(db, "clipboard.json"); err != nil {
		fmt.Printf("error exporting data to JSON: %v\n", err)
		return
	}

	fmt.Println("Data exported successfully.")
}

func exportToJSON(db *gorm.DB, filename string) error {
	var clips []ClipboardEntry
	if err := db.Find(&clips).Error; err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(clips); err != nil {
		return err
	}

	return nil
}
