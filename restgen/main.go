package main

import (
	"embed"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/pakohan/go/sqlrepo"
)

var (
	//go:embed sql
	sqlDir embed.FS
	repo   = sqlrepo.New(log.Default(), sqlDir, "sql")
)

func main() {
	log.SetFlags(log.Lshortfile)

	path, err := getConfigPath()
	if err != nil {
		log.Fatal(err.Error())
	}

	cfg, err := readConfig(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = os.Chdir(cfg.BaseDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	exists, err := checkFileExists(filepath.Join(cfg.BaseDir, "go.mod"))
	if err != nil {
		log.Fatal(err.Error())
	}

	if !exists {
		cmd := exec.Command("go", "mod", "init", cfg.ProjectImportPath)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	err = generateFiles(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}
