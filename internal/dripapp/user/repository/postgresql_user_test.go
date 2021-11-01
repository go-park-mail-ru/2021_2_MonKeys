package repository

import (
	"context"
	"database/sql"
	"dripapp/internal/dripapp/models"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// func TestPostgresqlUsers(t *testing.T, db *sql.DB) {
// 	conn, err := getDataBase()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	rep := PostgreUserRepo{
// 		Conn: *conn,
// 	}

// 	lu := models.LoginUser{
// 		Email:    "valid@valid.ru",
// 		Password: "!Nagdimaev2001",
// 	}
// 	lu2 := models.LoginUser{
// 		Email:    "valid@valid.ru",
// 		Password: "!Nagdimaev2001",
// 	}
// u := models.User{
// 	ID:       1,
// 	Email:    "valid@valid.ru",
// 	Password: "!Nagdimaev2001",
// }

// 	_ = true &&

// 		t.Run("create", func(t *testing.T) {
// 			user, err := rep.CreateUser(context.TODO(), lu)

// 			assert.Nil(t, err)
// 			assert.Equal(t, u, user)

// 			_, err = rep.CreateUser(context.TODO(), lu2)

// 			assert.NotNil(t, err)
// 		})
// }

// func addCreateUserSuport(mock sqlmock.Sqlmock) {
// 	rows := sqlmock.NewRows([]string{"id", "email", "password"}).AddRow(1, "valid@valid.ru", "!Nagdimaev2001")
// 	mock.ExpectQuery(CreateUserQuery).WithArgs("valid@valid.ru", "!Nagdimaev2001").WillReturnRows(rows)
// 	mock.ExpectQuery(CreateUserQuery).WithArgs("valid@valid.ru", "lol").WillReturnError(sql.ErrTxDone)
// }

// func addGetUserSupport(mock sqlmock.Sqlmock) {
// 	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description"}).
// 		AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", "всем привет")
// 	mock.ExpectQuery(GetUserQuery).WithArgs("valid@valid.ru").WillReturnRows(rows)
// 	mock.ExpectQuery(GetUserQuery).WithArgs("noexists").WillReturnError(sql.ErrNoRows)
// }

// func addGetImgsByIDSupport(mock sqlmock.Sqlmock) {
// 	rows := sqlmock.NewRows([]string{"imgs"}).AddRow(pq.Array([]string{"tag1", "tag2"}))
// 	mock.ExpectQuery(GetImgsByIDQuery).WithArgs(1).WillReturnRows(rows)
// 	mock.ExpectQuery(GetImgsByIDQuery).WithArgs(0).WillReturnError(sql.ErrNoRows)
// }

// func addGetTagsByIDSupport(mock sqlmock.Sqlmock) {
// 	rows := sqlmock.NewRows([]string{"tag_name"}).
// 		AddRow("tag1").
// 		AddRow("tag2")
// 	mock.ExpectQuery(GetTagsByIdQuery).WithArgs(1).WillReturnRows(rows)
// 	mock.ExpectQuery(GetTagsByIdQuery).WithArgs(0).WillReturnError(sql.ErrNoRows)
// }

func getDataBase() (*sqlx.DB, error) {
	db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	// addCreateUserSuport(mock)
	// addGetUserSupport(mock)
	// addGetImgsByIDSupport(mock)
	// addGetTagsByIDSupport(mock)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, err
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := &PostgreUserRepo{
		Conn: *sqlxDB,
	}

	u := models.User{
		ID:          1,
		Name:        "Ilyagu",
		Email:       "valid@valid.ru",
		Password:    "!Nagdimaev2001",
		Date:        "2001-06-29",
		Description: "всем привет",
		Imgs:        []string{"img1", "img2"},
	}

	// good query
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description", "imgs"}).
		AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", "всем привет", pq.Array([]string{"img1", "img2"}))
	mock.ExpectQuery("select id, name, email, password, date, description, imgs").
		WithArgs("valid@valid.ru").WillReturnRows(rows)

	user, err := repo.GetUser(context.TODO(), "valid@valid.ru")

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(user, u) {
		t.Errorf("results not match, want %v, have %v", u, user)
		return
	}

	// query error
	mock.ExpectQuery("select id, name, email, password, date, description, imgs").
		WithArgs("noexists").WillReturnError(sql.ErrNoRows)

	user, err = repo.GetUser(context.TODO(), "noexists")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// // row scan error
	// rows = sqlmock.NewRows([]string{"id", "title"}).
	// 	AddRow(1, "title")

	// mock.
	// 	ExpectQuery("SELECT id, title, updated, description FROM items WHERE").
	// 	WithArgs(elemID).
	// 	WillReturnRows(rows)

	// _, err = repo.GetByID(elemID)
	// if err := mock.ExpectationsWereMet(); err != nil {
	// 	t.Errorf("there were unfulfilled expectations: %s", err)
	// 	return
	// }
	// if err == nil {
	// 	t.Errorf("expected error, got nil")
	// 	return
	// }

}

// func TestCreateUser(t *testing.T) {
// 	lu := models.LoginUser{
// 		Email:    "valid@valid.ru",
// 		Password: "!Nagdimaev2001",
// 	}
// 	lu2 := models.LoginUser{
// 		Email:    "valid@valid.ru",
// 		Password: "!Nagdimaev2001",
// 	}
// 	u := models.User{
// 		ID:       1,
// 		Email:    "valid@valid.ru",
// 		Password: "!Nagdimaev2001",
// 	}
// conn, err := getDataBase()
// if err != nil {
// 	t.Error(err)
// }
// rep := PostgreUserRepo{
// 	Conn: *conn,
// }

// 	user, err := rep.CreateUser(context.TODO(), lu)

// 	assert.Nil(t, err)
// 	assert.Equal(t, u, user)

// 	_, err = rep.CreateUser(context.TODO(), lu2)

// 	assert.NotNil(t, err)
// }

// func TestGetUser(t *testing.T) {
// 	// u := models.User{
// 	// 	ID:          1,
// 	// 	Name:        "Ilyagu",
// 	// 	Email:       "valid@valid.ru",
// 	// 	Password:    "!Nagdimaev2001",
// 	// 	Date:        "2001-06-29",
// 	// 	Description: "всем привет",
// 	// 	// Imgs:        []string{"lol", "kek"},
// 	// }
// 	conn, err := getDataBase()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// rep := PostgreUserRepo{
// 	// 	Conn: *conn,
// 	// }

// 	// user, err := rep.GetUser(context.TODO(), "valid@valid.ru")
// 	// if err != nil {
// 	// 	t.Error(err)
// 	// }

// 	var RespUser models.User
// 	// err = conn.GetContext(context.TODO(), &RespUser, GetUserQuery, "valid@valid.ru")
// 	// err = conn.QueryRow(GetUserQuery, "valid@valid.ru").Scan(&RespUser.ID, &RespUser.Name, &RespUser.Email, &RespUser.Password,
// 	// 	&RespUser.Date, &RespUser.Description)
// 	// if err != nil {
// 	// 	t.Error(err)
// 	// }

// 	// err = conn.Select(&RespUser.Tags, GetTagsByIdQuery, 1)
// 	// if err != nil {
// 	// 	t.Error()
// 	// }

// 	if err = conn.QueryRow(GetImgsByIDQuery, 1).Scan(pq.Array(&RespUser.Imgs)); err != nil {
// 		t.Error(err)
// 	}

// 	assert.Nil(t, err)
// 	// assert.Equal(t, u, RespUser)
// }
