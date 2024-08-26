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
		CODE_ODESSEY_DATABASE_URI=postgres://user:password@localhost:5432/codeodessey_db?sslmode=disable
        CODE_ODESSEY_DB_HOST=localhost
        CODE_ODESSEY_DB_PASSWORD=password
        CODE_ODESSEY_DB_PORT=5432
        CODE_ODESSEY_DB_NAME=codeodessey_db
        CODE_ODESSEY_DB_SSL_MODE=disable
        CODE_ODESSEY_DB_USER=user`

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
	require.Equal(t, "localhost", config.DBHost)
	require.Equal(t, "password", config.DBPassword)
	require.Equal(t, "5432", config.DBPort)
	require.Equal(t, "codeodessey_db", config.DBName)
	require.Equal(t, "disable", config.DBSslMode)
	require.Equal(t, "user", config.DBUser)

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

//nolint:paralleltest
func TestLoadConfigEnvironmentVariablesOverride(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()
	// Create a temporary .env file
	envFilePath := filepath.Join(tmpDir, ".env")
	envFileContent := `CODE_ODESSEY_DB_DRIVER=postgres
		CODE_ODESSEY_DATABASE_URI=postgres://user:password@localhost:5432/codeodessey_db?sslmode=disable
        CODE_ODESSEY_DB_HOST=localhost
        CODE_ODESSEY_DB_PASSWORD=password
        CODE_ODESSEY_DB_PORT=5432
        CODE_ODESSEY_DB_NAME=codeodessey_db
        CODE_ODESSEY_DB_SSL_MODE=disable
        CODE_ODESSEY_DB_USER=user`

	err := os.WriteFile(envFilePath, []byte(envFileContent), 0o644)
	require.NoError(t, err)

	// Set environment variables to override the values in the .env file
	err = os.Setenv("CODE_ODESSEY_DB_DRIVER", "mysql")
	require.NoError(t, err)
	err = os.Setenv(
		"CODE_ODESSEY_DATABASE_URI",
		"mysql://user:password@localhost:3306/codeodessey_db",
	)
	require.NoError(t, err)
	err = os.Setenv("CODE_ODESSEY_DB_HOST", "newhost")
	require.NoError(t, err)
	err = os.Setenv("CODE_ODESSEY_DB_PASSWORD", "newpassword")
	require.NoError(t, err)
	err = os.Setenv("CODE_ODESSEY_DB_PORT", "3306")
	require.NoError(t, err)
	err = os.Setenv("CODE_ODESSEY_DB_NAME", "new_db")
	require.NoError(t, err)
	err = os.Setenv("CODE_ODESSEY_DB_SSL_MODE", "prefer")
	require.NoError(t, err)
	err = os.Setenv("CODE_ODESSEY_DB_USER", "newuser")
	require.NoError(t, err)

	defer func() {
		err = os.Unsetenv("CODE_ODESSEY_DB_DRIVER")
		require.NoError(t, err)
		err = os.Unsetenv("CODE_ODESSEY_DATABASE_URI")
		require.NoError(t, err)
		err = os.Unsetenv("CODE_ODESSEY_DB_HOST")
		require.NoError(t, err)
		err = os.Unsetenv("CODE_ODESSEY_DB_PASSWORD")
		require.NoError(t, err)
		err = os.Unsetenv("CODE_ODESSEY_DB_PORT")
		require.NoError(t, err)
		err = os.Unsetenv("CODE_ODESSEY_DB_NAME")
		require.NoError(t, err)
		err = os.Unsetenv("CODE_ODESSEY_DB_SSL_MODE")
		require.NoError(t, err)
		err = os.Unsetenv("CODE_ODESSEY_DB_USER")
		require.NoError(t, err)
	}()
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
	require.Equal(t, "newhost", config.DBHost)
	require.Equal(t, "newpassword", config.DBPassword)
	require.Equal(t, "3306", config.DBPort)
	require.Equal(t, "new_db", config.DBName)
	require.Equal(t, "prefer", config.DBSslMode)
	require.Equal(t, "newuser", config.DBUser)
}
