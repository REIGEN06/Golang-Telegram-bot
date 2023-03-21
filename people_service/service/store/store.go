package store

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	conn *sqlx.DB
}

type People struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// NewStore creates new database connection
func NewStore(connString string) *Store {
	conn, err := sqlx.Connect("pgx", connString)
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"../../migrations/1_initial.up.sql",
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	m.Up()

	return &Store{
		conn: conn,
	}
}

func (s *Store) ListPeople() ([]People, error) {
	people := make([]People, 0)
	var (
		name string
		id   int
	)
	//err := s.conn.SelectContext(ctx, &people, `SELECT name, address FROM people`)

	rows, err := s.conn.Query(`SELECT * FROM people`)
	if err != nil {
		return nil, fmt.Errorf("store ListPeople error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id, &name); err != nil {
			fmt.Fprintf(os.Stderr, "scan failed: %v\n", err)
			return nil, nil
		}

		people = append(people, People{ID: id, Name: name})
	}
	return people, err
}

func (s *Store) GetPeopleByID(id int) (People, error) {
	//err := s.conn.SelectContext(ctx, pp, `SELECT name, address FROM people WHERE People.ID = id`)
	var name string
	row := s.conn.QueryRow(`SELECT * FROM people WHERE id = ?`, id)

	if err := row.Scan(&id, &name); err != nil {
		return People{}, fmt.Errorf("store GetPeopleByID error: %w", err)
	}

	return People{ID: id, Name: name}, nil
}
