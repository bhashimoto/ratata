package front

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bhashimoto/ratata/types"
)

type WebAppConfig struct {
	Templates *template.Template
	BaseURL   string
	RootDir   string
	Client    *http.Client
}

func (cfg *WebAppConfig) Init(rootDir string, baseURL string) {
	cfg.Client = &http.Client{}
	cfg.BaseURL = baseURL
	cfg.RootDir = rootDir

	err := 	cfg.LoadTemplates()
	if err != nil {
		log.Fatal(err)
	}
}

func (cfg *WebAppConfig) RespondWithPageNotFound(w http.ResponseWriter) {
	cfg.RespondWithError(w, http.StatusNotFound, "Page not found")
	return
}

func (cfg *WebAppConfig) RespondWithError(w http.ResponseWriter, code int, msg string) {
	data := struct {
		Error string
	}{
		Error: msg,
	}

	w.WriteHeader(code)
	cfg.Templates.ExecuteTemplate(w, "error", data)
}


func (cfg *WebAppConfig) LoadTemplates() (error) {
	files, err := getTemplateFilesRecursive(cfg.RootDir)
	if err != nil {
		log.Println(err)
		return err
	}

	cfg.Templates, err = template.New("").ParseFiles(files...)
	return err
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
			if filepath.Ext(entry.Name()) == ".html" {
				log.Println("Loading template:", entry.Name())
				files = append(files, filepath.Join(path, entry.Name()))
			}
		}
	}
	return files, nil
}

func (cfg *WebAppConfig) sendRequest(endpoint string, method string, header *http.Header, body *bytes.Reader) (*http.Response, error) {
	fullPath := cfg.BaseURL + endpoint
	log.Println("Sending request to", fullPath)
	req, err := http.NewRequest(method, fullPath, body)
	if err != nil {
		log.Println("Error creating request:", err)
		return &http.Response{}, err
	}

	if header != nil {
		req.Header = *header
	}

	resp, err := cfg.Client.Do(req)
	if err != nil {
		log.Println("Error doing request at WebAppConfig.sendRequest:", err)
		return &http.Response{}, err
	}
	return resp, nil
}

func (cfg *WebAppConfig) fetchAccounts() ([]types.Account, error) {
	log.Println("fetching accounts")
	resp, err := cfg.sendRequest("accounts", "GET", nil, bytes.NewReader([]byte{}))
	if err != nil {
		log.Println("Error getting account")
		return []types.Account{}, err
	}

	accounts, err := cfg.responseToAccounts(resp)
	if err != nil {
		log.Println("error converting response")
		return []types.Account{}, err
	}
	log.Println("fetched successfully")
	return accounts, nil
}

func (cfg *WebAppConfig) fetchAccountBalance(accountId string) ([]types.Balance, []types.Payment, error) {
	resp, err := cfg.sendRequest(fmt.Sprintf("accounts/%s/balance", accountId), "GET", nil, bytes.NewReader([]byte{}))
	if err != nil {
		return []types.Balance{}, []types.Payment{}, err
	}
	defer resp.Body.Close()
	data := struct {
		Balances []types.Balance `json:"balances"`
		Payments []types.Payment `json:"payments"`
	}{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return []types.Balance{}, []types.Payment{}, err
	}
	return data.Balances, data.Payments, nil
}

func (cfg *WebAppConfig) fetchAccount(id string) (types.Account, error) {
	resp, err := cfg.sendRequest(fmt.Sprintf("accounts/%s", id), "GET", nil, bytes.NewReader([]byte{}))
	if err != nil {
		log.Println("error in the request in fetchAccount")
		return types.Account{}, err
	}
	defer resp.Body.Close()
	acc, err := cfg.responseToAccount(resp)
	if err != nil {
		log.Println("error converting response")
		return types.Account{}, err
	}
	return acc, nil

}

func (cfg *WebAppConfig) responseToAccount(resp *http.Response) (types.Account, error) {
	decoder := json.NewDecoder(resp.Body)
	acc := []types.Account{}
	err := decoder.Decode(&acc)
	if err != nil {
		return types.Account{}, err
	}

	return acc[0], nil
}

func (cfg *WebAppConfig) responseToAccounts(resp *http.Response) ([]types.Account, error) {
	decoder := json.NewDecoder(resp.Body)
	accs := []types.Account{}
	err := decoder.Decode(&accs)
	if err != nil {
		return []types.Account{}, err
	}
	return accs, nil
}
