package db

import (
	"connectfour/internal/model"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type MariaDbUserRepository struct {
	db *sql.DB
}

func NewMariaDbUserRepository() *MariaDbUserRepository {
	return &MariaDbUserRepository{
		db: connect(),
	}
}

func (r MariaDbUserRepository) Create(u model.User) (model.User, error) {
	result, err := r.db.Exec("INSERT INTO user (email, name, token) VALUES (?, ?, ?)", u.Email, u.Name, u.Token)
	if err != nil {
		log.Errorf("Error inserting user into the database: %v\n", err)
		return model.User{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Errorf("Error getting last insert ID: %v\n", err)
		return model.User{}, err
	}
	u.Id = id
	return u, nil
}

func (r MariaDbUserRepository) FindByEmail(email string) (model.User, error) {
	row := r.db.QueryRow("SELECT id, email, name, token FROM user WHERE email = ?", email)
	u := model.User{}
	if err := row.Scan(&u.Id, &u.Email, &u.Name, &u.Token); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debugf("Requested user '%s' not found (%v)\n", email, err)
			return model.User{}, nil
		} else {
			log.Errorf("Error getting user values from the database result: %v\n", err)
			return model.User{}, err
		}
	}
	return u, nil
}
