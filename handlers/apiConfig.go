package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bhashimoto/ratata/internal/database"
)

type ApiConfig struct {
	DB *database.Queries
	FS *http.Handler
	Templates map[string]*template.Template
	AccountCache map[string]*AccountData
}

func (cfg *ApiConfig) AddAccountDataToCache(id string, ad *AccountData) error {
	if cfg.AccountCache == nil {
		cfg.AccountCache = map[string]*AccountData{}
	}
	cfg.AccountCache[id] = ad
	return nil
}

func (cfg *ApiConfig) LoadTemplates(root string) (*template.Template, error) {
	files, err := getTemplateFilesRecursive(root)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, file := range files {
		name := filepath.Base(file)
		ts, err := template.New(name).ParseFiles(file)
		if err != nil {
			log.Fatal(err)
		}
		cfg.Templates[name] = ts
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return tmpl, nil
}

func (cfg *ApiConfig) CreateTemplteCache() error {
	cfg.Templates = make(map[string]*template.Template)
	return nil
}


func getTemplateFilesRecursive(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return []string{}, err
	}
	files := []string{}

	for _, entry := range entries {
		if entry.IsDir() {
			subPath := filepath.Join(path, entry.Name())
			subFiles, err := getTemplateFilesRecursive(subPath)
			if err != nil {
				return []string{}, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, filepath.Join(path, entry.Name()))
		}
	}
	return files, nil
}
