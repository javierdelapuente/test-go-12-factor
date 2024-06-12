package service

import (
	"database/sql"
	"errors"
	"log"

	"github.com/javierdelapuente/test-go-12-factor/config"
)

type Service struct {
	CharmConfig config.CharmConfig
}

func (s *Service) CheckMysqlStatus() (err error) {
	return errors.New("Not Implemented")
}

func (s *Service) CheckPostgresqlStatus() (err error) {
	// This can result in a pointer dereference
	db, err := sql.Open("pgx", *s.CharmConfig.Integrations.PostgresqlUrl)
	if err != nil {
		return
	}
	defer db.Close()

	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		return
	}
	log.Printf("postgresql version %s.", version)
	return
}
