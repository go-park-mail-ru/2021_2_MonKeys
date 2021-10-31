package repository_test

import (
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/hasher"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func addInsertSupport(mock sqlmock.Sqlmock) {
	hashedPass, err := hasher.HashAndSalt(nil, "!Nagdimaev2001")
	if err != nil {
		log.Fatal("ahahah")
	}
	query := `INSERT into profile(
		email,
		password)
		VALUES ($1,$2)
		RETURNING id, email, password;`

	rows := sqlmock.NewRows([]string{"id", "email", "password"}).AddRow(1, "valid@valid.ru", hashedPass)
	mock.ExpectQuery(query).WithArgs("valid@valid.ru", hashedPass).WillReturnRows(rows)
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
	hashedPass, err := hasher.HashAndSalt(nil, "!Nagdimaev2001")
	if err != nil {
		log.Fatal("hash error")
	}
	lu := models.LoginUser{
		Email:    "valid@valid.ru",
		Password: hashedPass,
	}
	// u := models.LoginUser{
	// 	ID:       1,
	// 	Email:    "valid@valid.ru",
	// 	Password: "!Nagdimaev2001",
	// }
	conn, err := getDataBase()
	if err != nil {
		t.Error(err)
	}
	// rep := repository.PostgreUserRepo{
	// 	Conn: *conn,
	// }

	var user models.User
	query := `INSERT into profile(
		email,
		password)
		VALUES ($1,$2)
		RETURNING id, email, password;`
	err = conn.QueryRow(query, lu.Email, lu.Password).Scan(&user)

	// user, err := rep.CreateUser(context.TODO(), lu)

	t.Error(user)
	// t.Error(u)
	assert.NotNil(t, err)
	// assert.Equal(t, u, user)
}
