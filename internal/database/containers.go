package database

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PGContainer struct {
	pgc *postgres.PostgresContainer
}

func NewPGContainer() (*PGContainer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbName := "rssagg"
	dbUser := "root"
	dbPassword := "secret"

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}
	ports, err := postgresContainer.Ports(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container ports: %w", err)
	}
	fmt.Println(ports)

	return &PGContainer{
		pgc: postgresContainer,
	}, nil
}

func (p *PGContainer) ConnString(ctx context.Context) (string, error) {
	return p.pgc.ConnectionString(ctx, "sslmode=disable")
}

func (p *PGContainer) Terminate(ctx context.Context) error {
	if err := p.pgc.Terminate(ctx); err != nil {
		return fmt.Errorf("failed to terminate container: %w", err)
	}

	return nil
}
