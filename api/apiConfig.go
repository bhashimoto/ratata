package api

import (
	"net/http"

	"github.com/bhashimoto/ratata/internal/database"
)

type ApiConfig struct {
	DB           *database.Queries
	FS           *http.Handler
	AccountCache map[string]*AccountData
}

func (cfg *ApiConfig) AddAccountDataToCache(id string, ad *AccountData) error {
	if cfg.AccountCache == nil {
		cfg.AccountCache = map[string]*AccountData{}
	}
	cfg.AccountCache[id] = ad
	return nil
}
