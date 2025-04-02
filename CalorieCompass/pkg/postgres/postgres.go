package postgres

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	defaultMaxPoolSize  = 1
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	DB *sqlx.DB
}

func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(pg)
	}

	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - sqlx.Connect: %w", err)
	}

	db.SetMaxOpenConns(pg.maxPoolSize)
	pg.DB = db

	return pg, nil
}

func (p *Postgres) Close() {
	if p.DB != nil {
		p.DB.Close()
	}
}

type Option func(*Postgres)

func MaxPoolSize(size int) Option {
	return func(pg *Postgres) {
		pg.maxPoolSize = size
	}
}