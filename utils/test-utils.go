package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/TheSandyDave/Media-Tags/domain"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func NewInMemoryDatabase(t *testing.T, models ...any) *gorm.DB {
	t.Helper()

	database, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		t.Error(err)
		return nil
	}
	database.Debug()
	if models == nil {
		require.NoError(t, database.AutoMigrate(domain.Models...))
	} else if len(models) > 0 {
		require.NoError(t, database.AutoMigrate(append(domain.Models, models...)...))
	}

	return database
}

func RetrieveResponse(t *testing.T, result any, code int, res *http.Response) bool {
	t.Helper()

	if code != res.StatusCode {
		t.Error("unexpected status code")
		return false
	}

	response, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("failed to read body of response")
		return false
	}

	if err := json.Unmarshal(response, &result); err != nil {
		t.Error("Failed unmarshalling result")
		return false
	}

	return true
}
