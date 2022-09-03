package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/g4web/otus_anti_brute_force/internal/config"
	"github.com/jmoiron/sqlx"
	// Postgres driver.
	_ "github.com/lib/pq"
)

const (
	NetworkTypeWhite = "white"
	NetworkTypeBlack = "black"
)

var ErrRowsAffected = errors.New("the number of affected rows is not equal to one")

type SQLStorage struct {
	db  *sqlx.DB
	ctx context.Context
}

func NewSQLStorage(ctx context.Context, c *config.Config) (*SQLStorage, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBUser, c.DBPassword, c.DBName)

	s := &SQLStorage{}
	err := s.connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *SQLStorage) connect(ctx context.Context, dsn string) error {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return err
	}

	s.db = db
	s.ctx = ctx

	return err
}

func (s *SQLStorage) AddToWhiteList(rawNetwork string) error {
	_, _, err := net.ParseCIDR(rawNetwork)
	if err != nil {
		return err
	}
	return s.insert(rawNetwork, NetworkTypeWhite)
}

func (s *SQLStorage) AddToBlackList(rawNetwork string) error {
	_, _, err := net.ParseCIDR(rawNetwork)
	if err != nil {
		return err
	}
	return s.insert(rawNetwork, NetworkTypeBlack)
}

func (s *SQLStorage) RemoveFromWhiteList(rawNetwork string) error {
	_, _, err := net.ParseCIDR(rawNetwork)
	if err != nil {
		return err
	}
	return s.delete(rawNetwork, NetworkTypeWhite)
}

func (s *SQLStorage) RemoveFromBlackList(rawNetwork string) error {
	_, _, err := net.ParseCIDR(rawNetwork)
	if err != nil {
		return err
	}
	return s.delete(rawNetwork, NetworkTypeBlack)
}

func (s *SQLStorage) GetWhiteLists() (map[string]*net.IPNet, error) {
	return s.selectRows(NetworkTypeWhite)
}

func (s *SQLStorage) GetBlackLists() (map[string]*net.IPNet, error) {
	return s.selectRows(NetworkTypeBlack)
}

func (s *SQLStorage) insert(network string, networkType string) error {
	query := `
				INSERT INTO network
					(network, type)
				VALUES
					($1, $2)
				;
	`

	_, err := s.db.ExecContext(
		s.ctx,
		query,
		network,
		networkType,
	)

	return err
}

func (s *SQLStorage) delete(network string, networkType string) error {
	query := `
				DELETE
				FROM
					network
				WHERE
					network = $1 AND type = $2
	`

	result, err := s.db.ExecContext(s.ctx, query, network, networkType)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return ErrRowsAffected
	}

	return nil
}

func (s *SQLStorage) selectRows(networkType string) (map[string]*net.IPNet, error) {
	networks := make(map[string]*net.IPNet)
	sqlQuery := `
		SELECT
		 network
		FROM
		  network
		WHERE
		  type = $1	
		;
	`
	rows, err := s.db.QueryxContext(s.ctx, sqlQuery, networkType)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var rawNetwork string

		err = rows.Scan(
			&rawNetwork,
		)

		if err != nil {
			return nil, err
		}

		_, network, err := net.ParseCIDR(rawNetwork)
		if err != nil {
			return nil, err
		}
		networks[rawNetwork] = network
	}

	return networks, nil
}
