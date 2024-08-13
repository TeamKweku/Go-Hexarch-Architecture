package postgres

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/teamkweku/code-odessey-hex-arch/config"
	"github.com/teamkweku/code-odessey-hex-arch/internal/adapters/outbound/postgres/sqlc"
)

const migrationsPath = "migrations"

//go:embed migrations/*.sql
var migrations embed.FS

type queries interface {
	sqlc.Querier
}

// client is a Posgres client
type Client struct {
	db      *pgxpool.Pool
	queries queries
}

// create a new PostgresClient
func New(ctx context.Context, url URL) (*Client, error) {
	config, err := pgxpool.ParseConfig(url.Expose())
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres at %q: %w", url, err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping db at %q: %w", url, err)
	}

	client := &Client{
		db:      db,
		queries: sqlc.New(db),
	}

	if err := client.migrate(); err != nil {
		db.Close()
		return nil, err
	}

	return client, nil
}

// close the database connection pool
func (c *Client) Close() error {
	c.db.Close()
	return nil
}

// migrate run all migrations
func (c *Client) migrate() error {
	migrator, err := newMigrator(c.db)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	defer func() {
		if _, closeErr := migrator.Close(); err != nil {
			fmt.Printf("error closing migrator: %v\n", closeErr)
		}
	}()

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate db: %w", err)
	}

	return nil
}

func newMigrator(
	pool *pgxpool.Pool,
) (*migrate.Migrate, error) {
	src, err := iofs.New(migrations, migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("create io.FS migration source: %w", err)
	}

	// Create a *sql.DB instance from the connection string
	connConfig, err := pgxpool.ParseConfig(pool.Config().ConnString())
	if err != nil {
		return nil, fmt.Errorf("parse connection string: %w", err)
	}

	db, err := sql.Open("pgx", connConfig.ConnString())
	if err != nil {
		return nil, fmt.Errorf("open sql.DB connection: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf(
			"create postgres migration driver from db: %w",
			err,
		)
	}
	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf(
			"create migrator from io.FS source and postgres driver: %w",
			err,
		)
	}

	return m, nil
}

// URL is a Postgres connection URL.
type URL struct {
	host     string
	port     string
	dbName   string
	user     string
	password string
	sslMode  string
}

func NewURL(cfg config.Config) URL {
	return URL{
		host:     cfg.DBHost,
		port:     cfg.DBPort,
		dbName:   cfg.DBName,
		user:     cfg.DBUser,
		password: cfg.DBPassword,
		sslMode:  cfg.DBSslMode,
	}
}

// GoString returns a Go-syntax representation of the URL, with the password
// redacted.
func (u URL) GoString() string {
	return fmt.Sprintf(
		"postgres.URL{host:%q, port:%q, dbName:%q, user:%q, password:REDACTED, sslMode:%q}",
		u.host,
		u.port,
		u.dbName,
		u.user,
		u.sslMode,
	)
}

// String returns a connection string with the password redacted.
func (u URL) String() string {
	return fmt.Sprintf(
		"postgres://%s:REDACTED@%s:%s/%s?sslmode=%s",
		u.user,
		u.host,
		u.port,
		u.dbName,
		u.sslMode,
	)
}

// Expose returns a connection string with the password exposed.
func (u URL) Expose() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		u.user,
		u.password,
		u.host,
		u.port,
		u.dbName,
		u.sslMode,
	)
}
