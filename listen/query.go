package listen

import (
	"fmt"
	"os"
	"text/template"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func QueryByCeationTime() {
	run("created_at")
}

func QueryByVideoTimestamp() {
	run("timestamp")
}

func ClipboardEntries(sortByField string) ([]ClipboardEntry, error) {
	db, err := gorm.Open(sqlite.Open("clipboard.db"), &gorm.Config{})
	if err != nil {
		return []ClipboardEntry{}, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.AutoMigrate(&ClipboardEntry{}); err != nil {
		return []ClipboardEntry{}, fmt.Errorf("error migrating database: %v", err)
	}

	var entries []ClipboardEntry
	if err := db.Find(&entries).Error; err != nil {
		return []ClipboardEntry{}, fmt.Errorf("error querying database: %v", err)
	}

	if err := db.Order(sortByField).Find(&entries).Error; err != nil {
		return []ClipboardEntry{}, fmt.Errorf("error sorting entries by %s: %v", sortByField, err)
	}

	return entries, nil
}

func ClipboardEntriesReverseByTimestamp() ([]ClipboardEntry, error) {
	db, err := gorm.Open(sqlite.Open("clipboard.db"), &gorm.Config{})
	if err != nil {
		return []ClipboardEntry{}, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.AutoMigrate(&ClipboardEntry{}); err != nil {
		return []ClipboardEntry{}, fmt.Errorf("error migrating database: %v", err)
	}

	var entries []ClipboardEntry
	if err := db.Find(&entries).Error; err != nil {
		return []ClipboardEntry{}, fmt.Errorf("error querying database: %v", err)
	}

	sortByField := "timestamp desc"

	if err := db.Order(sortByField).Find(&entries).Error; err != nil {
		return []ClipboardEntry{}, fmt.Errorf("error sorting entries by %s: %v", sortByField, err)
	}

	return entries, nil
}

func run(sortByField string) {
	db, err := gorm.Open(sqlite.Open("clipboard.db"), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		return
	}

	if err := db.AutoMigrate(&ClipboardEntry{}); err != nil {
		fmt.Printf("Error migrating database: %v\n", err)
		return
	}

	var entries []ClipboardEntry
	if err := db.Find(&entries).Error; err != nil {
		fmt.Printf("Error querying database: %v\n", err)
		return
	}

	if err := db.Order(sortByField).Find(&entries).Error; err != nil {
		fmt.Printf("Error sorting entries by %s: %v\n", sortByField, err)
		return
	}

	tmpl := `Clipboard Entries:
{{ range . }}* {{ .Content }}
{{ end }}`

	t := template.Must(template.New("clipboard").Parse(tmpl))
	if err := t.Execute(os.Stdout, entries); err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}
}
