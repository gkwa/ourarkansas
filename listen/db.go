package listen

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

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

func RunExportToSQL() error {
	cmd := exec.Command("sqlite3", "clipboard.db", ".dump")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	var data []byte
	buf := make([]byte, 1024)
	for {
		n, err := stdout.Read(buf)
		if err != nil {
			break
		}
		data = append(data, buf[:n]...)
	}

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Command exited with non-zero status: %v\n", exitErr)
			return nil
		}
		return fmt.Errorf("error waiting for command: %v", err)
	}

	if len(data) == 0 {
		fmt.Println("No data received on stdout, nothing written to clipboard.sql")
	} else {
		timestamp := time.Now().Unix()
		backupFileName := fmt.Sprintf("clipboard_%d.sql", timestamp)
		if err := os.WriteFile(backupFileName, data, 0o600); err != nil {
			return fmt.Errorf("error writing to backup file: %v", err)
		}
		fmt.Println("Data written to clipboard.sql")
	}

	return nil
}
