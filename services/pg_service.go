package services

import (
	"github.com/go-pg/pg/v9"
)

type PGService struct {
	DB *pg.DB
}

func (s *PGService) Connect(address string, user string, password string, database string, poolSize int) (err error) {
	db := pg.Connect(&pg.Options{
		Addr:     address,
		User:     user,
		Password: password,
		Database: database,
		PoolSize: poolSize,
	})

	// Test connection
	var n int
	_, err = db.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		return err
	}

	s.DB = db

	return nil
}

func (s *PGService) Kill() {
	s.DB.Close()
}
