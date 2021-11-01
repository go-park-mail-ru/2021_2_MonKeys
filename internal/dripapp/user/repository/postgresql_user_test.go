package repository

import (
	"context"
	"database/sql"
	"dripapp/internal/dripapp/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func testPostgresqlUsers(t *testing.T, db *sql.DB) {
	ctx := context.Background()
	conn, err := getDataBase()
	if err != nil {
		t.Error(err)
	}
	rep := PostgreUserRepo{
		Conn: *conn,
	}

	lu := models.LoginUser{
		Email:    "valid@valid.ru",
		Password: "!Nagdimaev2001",
	}
	// lu2 := models.LoginUser{
	// 	Email:    "valid@valid.ru",
	// 	Password: "!Nagdimaev2001",
	// }
	// u := models.User{
	// 	ID:       1,
	// 	Email:    "valid@valid.ru",
	// 	Password: "!Nagdimaev2001",
	// }

	_ = true &&

		t.Run("create", func(t *testing.T) {
			var err error
			_, err = rep.CreateUser(ctx, lu)
			assert.NoError(t, err)
		})
}

func addCreateUserSuport(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"id", "email", "password"}).AddRow(1, "valid@valid.ru", "!Nagdimaev2001")
	mock.ExpectQuery(CreateUserQuery).WithArgs("valid@valid.ru", "!Nagdimaev2001").WillReturnRows(rows)
	mock.ExpectQuery(CreateUserQuery).WithArgs("valid@valid.ru", "lol").WillReturnError(sql.ErrTxDone)
}

func addGetUserSupport(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description"}).
		AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", "всем привет")
	mock.ExpectQuery(GetUserQuery).WithArgs("valid@valid.ru").WillReturnRows(rows)
	mock.ExpectQuery(GetUserQuery).WithArgs("noexists").WillReturnError(sql.ErrNoRows)
}

func addGetImgsByIDSupport(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"imgs"}).AddRow(pq.Array([]string{"tag1", "tag2"}))
	mock.ExpectQuery(GetImgsByIDQuery).WithArgs(1).WillReturnRows(rows)
	mock.ExpectQuery(GetImgsByIDQuery).WithArgs(0).WillReturnError(sql.ErrNoRows)
}

func addGetTagsByIDSupport(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"tag_name"}).
		AddRow("tag1").
		AddRow("tag2")
	mock.ExpectQuery(GetImgsByIDQuery).WithArgs(1).WillReturnRows(rows)
	mock.ExpectQuery(GetImgsByIDQuery).WithArgs(0).WillReturnError(sql.ErrNoRows)
}

func getDataBase() (*sqlx.DB, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	addCreateUserSuport(mock)
	addGetUserSupport(mock)
	addGetImgsByIDSupport(mock)
	addGetTagsByIDSupport(mock)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, err
}

func TestCreateUser(t *testing.T) {
	lu := models.LoginUser{
		Email:    "valid@valid.ru",
		Password: "!Nagdimaev2001",
	}
	lu2 := models.LoginUser{
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
	rep := PostgreUserRepo{
		Conn: *conn,
	}

	user, err := rep.CreateUser(context.TODO(), lu)

	assert.Nil(t, err)
	assert.Equal(t, u, user)

	_, err = rep.CreateUser(context.TODO(), lu2)

	assert.NotNil(t, err)
}

func TestGetUser(t *testing.T) {
	u := models.User{
		ID:          1,
		Name:        "Ilyagu",
		Email:       "valid@valid.ru",
		Password:    "!Nagdimaev2001",
		Date:        "2001-06-29",
		Description: "всем привет",
		// Imgs:        []string{"lol", "kek"},
	}
	conn, err := getDataBase()
	if err != nil {
		t.Error(err)
	}
	rep := PostgreUserRepo{
		Conn: *conn,
	}

	user, err := rep.GetUser(context.TODO(), "valid@valid.ru")

	assert.Nil(t, err)
	assert.Equal(t, u, user)
}
