//go:build mage
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const dbName = "words.db"

type Word struct {
	Japanese string `json:"japanese"`
	Romaji   string `json:"romaji"`
	English  string `json:"english"`
}

// InitDB initializes the SQLite database
func InitDB() error {
	if _, err := os.Stat(dbName); err == nil {
		fmt.Printf("Database %s already exists\n", dbName)
		return nil
	}

	file, err := os.Create(dbName)
	if err != nil {
		return fmt.Errorf("error creating database: %v", err)
	}
	file.Close()

	fmt.Printf("Created database %s\n", dbName)
	return nil
}

// Migrate runs all migration files in order
func Migrate() error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	files, err := ioutil.ReadDir("db/migrations")
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %v", err)
	}

	var migrations []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sql") {
			migrations = append(migrations, f.Name())
		}
	}
	sort.Strings(migrations)

	for _, migration := range migrations {
		fmt.Printf("Running migration: %s\n", migration)
		content, err := ioutil.ReadFile(filepath.Join("db/migrations", migration))
		if err != nil {
			return fmt.Errorf("error reading migration file %s: %v", migration, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("error executing migration %s: %v", migration, err)
		}
	}

	return nil
}

// Seed imports data from JSON files in the seeds directory
func Seed() error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	files, err := ioutil.ReadDir("db/seeds")
	if err != nil {
		return fmt.Errorf("error reading seeds directory: %v", err)
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".json") {
			continue
		}

		groupName := strings.TrimSuffix(f.Name(), ".json")
		fmt.Printf("Seeding group: %s\n", groupName)

		// Create group
		result, err := db.Exec("INSERT INTO groups (name) VALUES (?)", groupName)
		if err != nil {
			return fmt.Errorf("error creating group %s: %v", groupName, err)
		}

		groupID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("error getting group ID: %v", err)
		}

		// Read and parse words
		content, err := ioutil.ReadFile(filepath.Join("db/seeds", f.Name()))
		if err != nil {
			return fmt.Errorf("error reading seed file %s: %v", f.Name(), err)
		}

		var words []Word
		if err := json.Unmarshal(content, &words); err != nil {
			return fmt.Errorf("error parsing seed file %s: %v", f.Name(), err)
		}

		// Insert words and create word-group associations
		for _, word := range words {
			result, err := db.Exec(
				"INSERT INTO words (japanese, romaji, english) VALUES (?, ?, ?)",
				word.Japanese, word.Romaji, word.English,
			)
			if err != nil {
				return fmt.Errorf("error inserting word %v: %v", word, err)
			}

			wordID, err := result.LastInsertId()
			if err != nil {
				return fmt.Errorf("error getting word ID: %v", err)
			}

			_, err = db.Exec(
				"INSERT INTO words_groups (word_id, group_id) VALUES (?, ?)",
				wordID, groupID,
			)
			if err != nil {
				return fmt.Errorf("error creating word-group association: %v", err)
			}
		}
	}

	return nil
}
