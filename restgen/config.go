package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
)

const configFile = "restgen.json"

type Config struct {
	ProjectImportPath string `json:"project_import_path"`

	DatabaseURL string  `json:"database_url"`
	TableSchema string  `json:"table_schema"`
	Models      []Model `json:"models"`
}

type Model struct {
	ProjectImportPath string   `json:"-"`
	TableName         string   `json:"table_name" db:"table_name"`
	FilterColumns     []string `json:"filter_columns"`
}

func makeConfig() error {
	exists, err := checkFileExists(configFile)
	if err != nil {
		return err
	} else if exists {
		return nil
	}

	cfg := &Config{}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Please enter postgres connection URL: ")
	scanner.Scan()
	cfg.DatabaseURL = scanner.Text()

	db, err := sqlx.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	fmt.Print("Please enter schema: ")
	scanner.Scan()
	cfg.TableSchema = strings.TrimSpace(scanner.Text())

	cfg.Models, err = getTables(db, repo, cfg.TableSchema)
	if err != nil {
		return err
	}

	fmt.Print("Please enter base import path for project: ")
	scanner.Scan()
	cfg.ProjectImportPath = strings.TrimSpace(scanner.Text())

	log.Printf("Creating json config file '%s'\n", configFile)
	f, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(cfg)
}

func readConfig() (*Config, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &Config{}
	err = json.NewDecoder(f).Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func checkFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, fmt.Errorf("err checking if file exists '%s': %s", path, err.Error())
	}
}
