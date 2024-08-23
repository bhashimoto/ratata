package handlerstest_tes

import (
	"reflect"
	"testing"
	"time"

	"github.com/bhashimoto/ratata/handlers"
	"github.com/bhashimoto/ratata/internal/database"
)

func TestDBUserToUser(t *testing.T) {
	cfg := handlers.ApiConfig{}
	dbUser := database.User{
		ID:         1,
		Name:       "Test user",
		CreatedAt:  1723121932,
		ModifiedAt: 1723121932,
	}

	got := cfg.DBUserToUser(dbUser)
	want := handlers.User{
		ID:         1,
		Name:       "Test user",
		CreatedAt:  time.Unix(1723121932, 0),
		ModifiedAt: time.Unix(1723121932, 0),
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("DBUserToUser: resulting structs are not equal")
	}

}
