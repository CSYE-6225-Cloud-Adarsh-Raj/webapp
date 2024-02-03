package setup

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDBConfig(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "pass")
	os.Setenv("DB_NAME", "dbname")

	expectedDSN := "host=localhost user=user dbname=dbname sslmode=disable password=pass"
	config := GetDBConfig()

	assert.Equal(t, expectedDSN, config.DSN)

	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
}
