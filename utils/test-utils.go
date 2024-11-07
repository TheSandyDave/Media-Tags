package utils

import (
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
