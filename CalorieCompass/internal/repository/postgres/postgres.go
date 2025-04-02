package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// PostgresDB - PostgreSQL veritabanı yapısı
type PostgresDB struct {
	db *sqlx.DB
}

// New - yeni PostgreSQL bağlantısı oluşturur
func New(db *sqlx.DB) *PostgresDB {
	return &PostgresDB{
		db: db,
	}
}

// GetDB - sqlx.DB nesnesini döndürür
func (p *PostgresDB) GetDB() *sqlx.DB {
	return p.db
}

// Close - veritabanı bağlantısını kapatır
func (p *PostgresDB) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return fmt.Errorf("db connection is nil")
}