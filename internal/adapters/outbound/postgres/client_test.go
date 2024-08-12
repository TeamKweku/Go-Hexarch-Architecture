package postgres

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/teamkweku/code-odessey-hex-arch/config"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("DB connection and migrations", func(t *testing.T) {
		t.Parallel()

		cfg, err := config.LoadConfig("../../../..")
		require.NoError(t, err)

		client, err := New(context.Background(), NewURL(cfg))
		assert.NoError(t, err)
		assert.NotNil(t, client)

		// Type assertion to access the underlying PostgresClient
		postgresClient, ok := client.(*PostgresClient)
		require.True(t, ok, "Client is not of type *PostgresClient")

		assert.NotNil(t, postgresClient.db)
		assert.NotNil(t, postgresClient.Queries)

		err = postgresClient.db.Ping(context.Background())
		assert.NoError(t, err)

		expectedVersion := latestMigrationVersion(t)
		migrator, err := newMigrator(postgresClient.db)
		require.NoError(t, err)

		gotVersion, dirty, err := migrator.Version()
		require.NoError(t, err)
		assert.Equal(t, expectedVersion, gotVersion)
		assert.False(t, dirty, "Latest migration is dirty")

		_ = postgresClient.Close()
	})

	t.Run("returns an error", func(t *testing.T) {
		t.Parallel()

		invalidCfg := config.Config{
			DBHost:     "invalid-host",
			DBPort:     "5432",
			DBName:     "invalid-db",
			DBUser:     "invliad-user",
			DBPassword: "invalid-password",
			DBSslMode:  "disable",
		}

		client, err := New(context.Background(), NewURL(invalidCfg))
		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

func Test_URL_New(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBName:     "codeodessey",
		DBUser:     "user",
		DBPassword: "password",
		DBSslMode:  "disable",
	}

	wantURL := URL{
		host:     cfg.DBHost,
		port:     cfg.DBPort,
		dbName:   cfg.DBName,
		user:     cfg.DBUser,
		password: cfg.DBPassword,
		sslMode:  cfg.DBSslMode,
	}

	assert.Equal(t, wantURL, NewURL(cfg))
}

func Test_URL_GoString(t *testing.T) {
	t.Parallel()

	url := URL{
		host:     "localhost",
		port:     "5432",
		dbName:   "realworld",
		user:     "user",
		password: "password",
		sslMode:  "disable",
	}

	want := `postgres.URL{host:"localhost", port:"5432", dbName:"realworld", user:"user", password:REDACTED, sslMode:"disable"}`
	assert.Equal(t, want, url.GoString())
}

func Test_URL_String(t *testing.T) {
	t.Parallel()

	url := URL{
		host:     "localhost",
		port:     "5432",
		dbName:   "realworld",
		user:     "user",
		password: "password",
		sslMode:  "disable",
	}

	want := "postgres://user:REDACTED@localhost:5432/realworld?sslmode=disable"
	assert.Equal(t, want, url.String())
}

func Test_URL_Expose(t *testing.T) {
	t.Parallel()

	url := URL{
		host:     "localhost",
		port:     "5432",
		dbName:   "realworld",
		user:     "user",
		password: "password",
		sslMode:  "disable",
	}

	want := "postgres://user:password@localhost:5432/realworld?sslmode=disable"
	assert.Equal(t, want, url.Expose())
}

func latestMigrationVersion(t *testing.T) uint {
	t.Helper()

	migrations, err := os.ReadDir(migrationsPath)
	require.NoError(t, err)

	var latestVersion uint
	for _, file := range migrations {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			versionStr := strings.Split(file.Name(), "_")[0]
			version, err := strconv.ParseUint(versionStr, 10, 32)
			require.NoError(t, err)
			if uint(version) > latestVersion {
				latestVersion = uint(version)
			}
		}
	}

	require.NotZero(t, latestVersion, "No valid migration files found")
	return latestVersion
}
