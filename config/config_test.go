package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfigFromEnv(t *testing.T) {
	t.Parallel()
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create a temporary .env file
	envFilePath := filepath.Join(tmpDir, ".env")
	envFileContent := `CODE_ODESSEY_DB_DRIVER=postgres
		CODE_ODESSEY_DATABASE_URI=postgres://user:password@localhost:5432/codeodessey_db?sslmode=disable`

	err := os.WriteFile(envFilePath, []byte(envFileContent), 0o644)
	require.NoError(t, err)

	// Load the config
	config, err := LoadConfig(tmpDir)
	require.NoError(t, err)

	// Validate the loaded config
	require.Equal(t, "postgres", config.DBDriver)
	require.Equal(
		t,
		"postgres://user:password@localhost:5432/codeodessey_db?sslmode=disable",
		config.DBSource,
	)

	// Set environment variables to override the values in the .env file
	err = os.Setenv("CODE_ODESSEY_DB_DRIVER", "mysql")
	require.NoError(t, err)

	err = os.Setenv(
		"CODE_ODESSEY_DATABASE_URI",
		"mysql://user:password@localhost:3306/codeodessey_db",
	)
	require.NoError(t, err)

	defer func() {
		err := os.Unsetenv("CODE_ODESSEY_DB_DRIVER")
		require.NoError(t, err)

		err = os.Unsetenv("CODE_ODESSEY_DATABASE_URI")
		require.NoError(t, err)
	}()
}

func TestLoadConfigEnvironmentVariablesOverride(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()
	// Create a temporary .env file
	envFilePath := filepath.Join(tmpDir, ".env")
	envFileContent := `CODE_ODESSEY_DB_DRIVER=postgres
		CODE_ODESSEY_DATABASE_URI=postgres://user:password@localhost:5432/codeodessey_db?sslmode=disable`
	err := os.WriteFile(envFilePath, []byte(envFileContent), 0o644)
	require.NoError(t, err)

	// Set environment variables to override the values in the .env file
	os.Setenv("CODE_ODESSEY_DB_DRIVER", "mysql")
	os.Setenv(
		"CODE_ODESSEY_DATABASE_URI",
		"mysql://user:password@localhost:3306/codeodessey_db",
	)
	defer os.Unsetenv("CODE_ODESSEY_DB_DRIVER")
	defer os.Unsetenv("CODE_ODESSEY_DATABASE_URI")

	// Load the config
	config, err := LoadConfig(tmpDir)
	require.NoError(t, err)

	// Validate the loaded config
	require.Equal(t, "mysql", config.DBDriver)
	require.Equal(
		t,
		"mysql://user:password@localhost:3306/codeodessey_db",
		config.DBSource,
	)
}
