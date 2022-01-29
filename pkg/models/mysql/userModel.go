package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/JeremyCoco/kubernetes-app/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (um *UserModel) Insert(username string, password string) error {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	_, err = um.DB.Exec(insertUser, username, string(hashed_password))

	if err != nil {
		var mySQLErr *mysql.MySQLError

		if errors.As(err, &mySQLErr) {
			if mySQLErr.Number == 1062 && strings.Contains(mySQLErr.Message, "username") {
				return models.ErrDuplicatedUsername
			}
		}

		return err
	}

	return nil
}

func (um *UserModel) Authenticate(username string, password string) (int, error) {
	var id int
	var hashed_password []byte

	err := um.DB.QueryRow(selectUserForAuth, username).Scan(&id, &hashed_password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		}

		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashed_password, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		}

		return 0, err
	}

	return id, nil
}

func (um *UserModel) GetById(id int) (*models.User, error) {
	user := &models.User{}

	err := um.DB.QueryRow(selectUserById, id).Scan(&user.Id, &user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecords
		}

		return nil, err
	}

	return user, nil
}
