package handlers

import (
	"net/http"

	"github.com/bhashimoto/ratata/internal/database"
)

type ApiConfig struct {
	DB *database.Queries
	FS *http.Handler
}
