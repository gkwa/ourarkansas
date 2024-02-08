package test1

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/glebarez/sqlite"
	"github.com/taylormonacelli/ourarkansas/youtube"
	"gorm.io/gorm"
)

func RunTest1() {
	startTime := time.Now()

	db, err := gorm.Open(sqlite.Open("clipboard.db"), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		return
	}

	err = db.AutoMigrate(&ClipboardEntry{})
	if err != nil {
		fmt.Printf("Error migrating database: %v\n", err)
		return
	}

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	for {
		select {
		case <-sig:
			fmt.Println("Received SIGINT. Exiting...")
			return

		case <-ticker.C:
			elapsed := time.Since(startTime)

			content, err := clipboard.ReadAll()
			if err != nil {
				fmt.Printf("[%s] Error reading clipboard: %v\n", elapsed, err)
				continue
			}

			content = strings.TrimSpace(content)

			if !isValidContent(content) {
				fmt.Printf("[%s] Content '%s' does not match the pattern. Skipping...\n", elapsed.Truncate(time.Second), content)
				continue
			}

			var count int64
			db.Model(&ClipboardEntry{}).Where("content = ?", content).Count(&count)
			if count > 0 {
				fmt.Printf("[%s] Content '%s' already exists in the database. Skipping...\n", elapsed.Truncate(time.Second), content)
				continue
			}

			result, err := youtube.DeconstructYouTubeURL(content)
			if err != nil {
				fmt.Printf("[%s] Error deconstructing YouTube URL: %v\n", elapsed.Truncate(time.Second), err)
				continue
			}

			entry := ClipboardEntry{
				Content:   content,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Timestamp: result.TimestampSeconds,
				VideoID:   result.VideoID,
			}
			if err := db.Create(&entry).Error; err != nil {
				fmt.Printf("Error inserting entry into database: %v\n", err)
				continue
			}

			fmt.Printf("Content '%s' written to database\n", content)
		}
	}
}

func isValidContent(content string) bool {
	return youtube.IsYoutubeURL(content)
}
