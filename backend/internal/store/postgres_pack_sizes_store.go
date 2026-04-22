package store

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	dbsqlc "pack_wise/internal/db/sqlc"
)

var ErrInvalidPackSizes = errors.New("invalid pack sizes")

type PostgresPackSizesStore struct {
	pool    *pgxpool.Pool
	queries *dbsqlc.Queries
}

func NewPostgresPackSizesStore(ctx context.Context, databaseURL string) (*PostgresPackSizesStore, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	return &PostgresPackSizesStore{
		pool:    pool,
		queries: dbsqlc.New(pool),
	}, nil
}

func (s *PostgresPackSizesStore) Close() {
	s.pool.Close()
}

func (s *PostgresPackSizesStore) GetPackSizes(ctx context.Context) ([]int, error) {
	packSizes, err := s.queries.ListPackSizes(ctx)
	if err != nil {
		return nil, fmt.Errorf("list pack sizes: %w", err)
	}

	sizes := make([]int, len(packSizes))
	for i, packSize := range packSizes {
		sizes[i] = int(packSize)
	}

	return sizes, nil
}

func (s *PostgresPackSizesStore) ReplacePackSizes(ctx context.Context, sizes []int) error {
	normalizedSizes, err := normalizePackSizes(sizes)
	if err != nil {
		return err
	}

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin replace pack sizes transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	queries := s.queries.WithTx(tx)
	if err := queries.DeletePackSizes(ctx); err != nil {
		return fmt.Errorf("delete pack sizes: %w", err)
	}

	for _, size := range normalizedSizes {
		if err := queries.CreatePackSize(ctx, int32(size)); err != nil {
			return fmt.Errorf("create pack size %d: %w", size, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit replace pack sizes transaction: %w", err)
	}

	return nil
}

func (s *PostgresPackSizesStore) DeletePackSize(ctx context.Context, size int) error {
	if size <= 0 {
		return ErrInvalidPackSizes
	}

	if err := s.queries.DeletePackSize(ctx, int32(size)); err != nil {
		return fmt.Errorf("delete pack size %d: %w", size, err)
	}

	return nil
}

func normalizePackSizes(sizes []int) ([]int, error) {
	if len(sizes) == 0 {
		return nil, ErrInvalidPackSizes
	}

	normalizedSizes := append([]int(nil), sizes...)
	sort.Ints(normalizedSizes)

	for i, size := range normalizedSizes {
		if size <= 0 {
			return nil, ErrInvalidPackSizes
		}

		if i > 0 && normalizedSizes[i-1] == size {
			return nil, ErrInvalidPackSizes
		}
	}

	return normalizedSizes, nil
}
