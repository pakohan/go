package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
)

const configFile = "restgen.json"

type Config struct {
	BaseDir           string `json:"base_dir"`
	ProjectImportPath string `json:"project_import_path"`

	DatabaseURL string  `json:"database_url"`
	TableSchema string  `json:"table_schema"`
	Models      []Model `json:"models"`
}

type Model struct {
	ProjectImportPath string `json:"-"`
	TableName         string `json:"table_name" db:"table_name"`
}

func getConfigPath() (string, error) {
	if len(os.Args) > 1 {
		fmt.Printf("checking for config file at '%s'\n", configFile)
		exists, err := checkFileExists(configFile)
		if err != nil {
			return "", err
		} else if exists {
			return configFile, nil
		}
	}

	scanner := bufio.NewScanner(os.Stdin)

	var baseDir string
	if len(os.Args) > 1 {
		baseDir = os.Args[1]
	} else {
		fmt.Printf("Please enter folder for your project: ")
		scanner.Scan()
		var err error
		baseDir, err = filepath.Abs(scanner.Text())
		if err != nil {
			return "", err
		}
	}

	p := filepath.Join(baseDir, configFile)
	fmt.Printf("checking for config file at '%s'\n", p)
	exists, err := checkFileExists(p)
	if err != nil {
		return "", err
	} else if exists {
		return p, nil
	}

	cfg := &Config{
		BaseDir: baseDir,
	}

	fmt.Print("Please enter postgres connection URL: ")
	scanner.Scan()
	cfg.DatabaseURL = scanner.Text()

	db, err := sqlx.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return "", err
	}
	defer db.Close()

	fmt.Print("Please enter schema: ")
	scanner.Scan()
	cfg.TableSchema = strings.TrimSpace(scanner.Text())

	cfg.Models, err = getTables(db, repo, cfg.TableSchema)
	if err != nil {
		return "", err
	}

	fmt.Print("Please enter base import path for project: ")
	scanner.Scan()
	cfg.ProjectImportPath = strings.TrimSpace(scanner.Text())

	p = filepath.Join(baseDir, configFile)
	log.Printf("Creating json config file '%s'\n", p)
	err = os.MkdirAll(baseDir, fs.ModePerm)
	if err != nil {
		return "", err
	}

	f, err := os.Create(p)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return p, json.NewEncoder(f).Encode(cfg)
}

func readConfig(path string) (*Config, error) {
	f, err := os.Open(path)
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
