package repository_test

import (
	"context"
	"dripapp/internal/dripapp/models"
	"dripapp/internal/dripapp/user/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func addInsertSupport(mock sqlmock.Sqlmock) {
	query := `INSERT into profile(
		email,
		password)
		VALUES ($1,$2)
		RETURNING id, email, password;`

	rows := sqlmock.NewRows([]string{"id", "email", "password"}).AddRow(1, "valid@valid.ru", "!Nagdimaev2001")
	mock.ExpectQuery(query).WithArgs("valid@valid.ru", "!Nagdimaev2001").WillReturnRows(rows)
}

func getDataBase() (*sqlx.DB, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	addInsertSupport(mock)
	// addGetByIdSupport(mock)
	// addAddSupport(mock)
	// addUpdateSupport(mock)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, err
}

func TestCreateUser(t *testing.T) {
	lu := models.LoginUser{
		Email:    "valid@valid.ru",
		Password: "!Nagdimaev2001",
	}
	u := models.User{
		ID:       1,
		Email:    "valid@valid.ru",
		Password: "!Nagdimaev2001",
	}
	conn, err := getDataBase()
	if err != nil {
		t.Error(err)
	}
	rep := repository.PostgreUserRepo{
		Conn: *conn,
	}

	var user models.User

	user, err = rep.CreateUser(context.TODO(), lu)

	assert.Nil(t, err)
	assert.Equal(t, u, user)
}
