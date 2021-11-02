package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"dripapp/internal/dripapp/models"
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func TestGetUser(t *testing.T) {
	t.Parallel()
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
		Tags:        []string{"tag1", "tag2"},
	}

	t.Run("good get user", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description", "imgs"}).
			AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", "всем привет", pq.Array([]string{"img1", "img2"}))
		mock.ExpectQuery("select id, name, email, password, date, description, imgs").
			WithArgs("valid@valid.ru").WillReturnRows(rows)

		rowsTags := sqlmock.NewRows([]string{"tag_name"}).
			AddRow("tag1").
			AddRow("tag2")
		mock.ExpectQuery("select tag_name").WithArgs(1).WillReturnRows(rowsTags)

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
	})
	t.Run("no tags", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description", "imgs"}).
			AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", "всем привет", pq.Array([]string{"img1", "img2"}))
		mock.ExpectQuery("select id, name, email, password, date, description, imgs").
			WithArgs("valid@valid.ru").WillReturnRows(rows)

		mock.ExpectQuery("select tag_name").WithArgs(1).WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUser(context.TODO(), "valid@valid.ru")

		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		u.Tags = []string(nil)
		if !reflect.DeepEqual(user, u) {
			t.Errorf("results not match, want \n%v\n, have \n%v\n", u, user)
			return
		}
	})
	t.Run("no users", func(t *testing.T) {
		mock.ExpectQuery("select id, name, email, password, date, description, imgs").
			WithArgs("noexists").WillReturnError(sql.ErrNoRows)

		_, err = repo.GetUser(context.TODO(), "noexists")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if err == nil {
			t.Errorf("expected error, got nil")
			return
		}
	})
	t.Run("error tags", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description", "imgs"}).
			AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", "всем привет", pq.Array([]string{"img1", "img2"}))
		mock.ExpectQuery("select id, name, email, password, date, description, imgs").
			WithArgs("valid@valid.ru").WillReturnRows(rows)

		mock.ExpectQuery("select tag_name").WithArgs(1).WillReturnError(fmt.Errorf("db_error"))

		_, err := repo.GetUser(context.TODO(), "valid@valid.ru")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if err == nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
	})
}

func TestCreateUser(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := &PostgreUserRepo{
		Conn: *sqlxDB,
	}

	lu := models.LoginUser{
		Email:    "valid@valid.ru",
		Password: "!Nagdimaev2001",
	}
	u := models.User{
		ID:       1,
		Email:    "valid@valid.ru",
		Password: "!Nagdimaev2001",
	}

	t.Run("good create user", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(1, "valid@valid.ru", "!Nagdimaev2001")
		mock.ExpectQuery("INSERT into").
			WithArgs("valid@valid.ru", "!Nagdimaev2001").WillReturnRows(rows)

		user, err := repo.CreateUser(context.TODO(), lu)

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
	})
	t.Run("error create user", func(t *testing.T) {
		mock.ExpectQuery("INSERT into").
			WithArgs("valid@valid.ru", "!Nagdimaev2001").WillReturnError(sql.ErrNoRows)

		_, err := repo.CreateUser(context.TODO(), lu)

		if err == nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	})
}

func TestGetUserByID(t *testing.T) {
	t.Parallel()
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
		Tags:        []string{"tag1", "tag2"},
	}

	t.Run("good get user by id", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description", "imgs"}).
			AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", "всем привет", pq.Array([]string{"img1", "img2"}))
		mock.ExpectQuery("select id, name, email, password, date, description, imgs").
			WithArgs(1).WillReturnRows(rows)

		rowsTags := sqlmock.NewRows([]string{"tag_name"}).
			AddRow("tag1").
			AddRow("tag2")
		mock.ExpectQuery("select tag_name").WithArgs(1).WillReturnRows(rowsTags)

		user, err := repo.GetUserByID(context.TODO(), 1)

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
	})
	t.Run("no tags", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description", "imgs"}).
			AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", "всем привет", pq.Array([]string{"img1", "img2"}))
		mock.ExpectQuery("select id, name, email, password, date, description, imgs").
			WithArgs(1).WillReturnRows(rows)

		mock.ExpectQuery("select tag_name").WithArgs(1).WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByID(context.TODO(), 1)

		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		u.Tags = []string(nil)
		if !reflect.DeepEqual(user, u) {
			t.Errorf("results not match, want \n%v\n, have \n%v\n", u, user)
			return
		}
	})
	t.Run("no users", func(t *testing.T) {
		mock.ExpectQuery("select id, name, email, password, date, description, imgs").
			WithArgs(0).WillReturnError(sql.ErrNoRows)

		_, err = repo.GetUserByID(context.TODO(), 0)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if err == nil {
			t.Errorf("expected error, got nil")
			return
		}
	})
	t.Run("error tags", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description", "imgs"}).
			AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", "всем привет", pq.Array([]string{"img1", "img2"}))
		mock.ExpectQuery("select id, name, email, password, date, description, imgs").
			WithArgs(1).WillReturnRows(rows)

		mock.ExpectQuery("select tag_name").WithArgs(1).WillReturnError(fmt.Errorf("db_error"))

		_, err := repo.GetUserByID(context.TODO(), 1)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if err == nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
	})
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
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
		Tags:        []string{"tag1", "tag2"},
	}

	// t.Run("error delete tags", func(t *testing.T) {
	// 	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description", "imgs"}).
	// 		AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", "всем привет", pq.Array([]string{"img1", "img2"}))
	// 	mock.ExpectQuery("update").WithArgs(
	// 		"Ilyagu",
	// 		"valid@valid.ru",
	// 		"2001-06-29",
	// 		"всем привет",
	// 		pq.Array([]string{"img1", "img2"}),
	// 	).WillReturnRows(rows)

	// 	// vals := []driver.Value{1, "anime", "music"}

	// 	mock.ExpectQuery("delete").WithArgs(1).WillReturnError(sql.ErrNoRows)
	// 	// mock.ExpectQuery("insert into profile_tag").WithArgs(nil).WillReturnError(nil)

	// 	// rowsTags := sqlmock.NewRows([]string{"tag_name"}).AddRow("tag1").AddRow("tag2")
	// 	// mock.ExpectQuery("select tag_name").WithArgs(1).WillReturnRows(rowsTags)

	// 	_, err := repo.UpdateUser(context.TODO(), u)

	// 	if err == nil {
	// 		t.Errorf("unexpected err: %s", err)
	// 		return
	// 	}
	// 	if err := mock.ExpectationsWereMet(); err != nil {
	// 		t.Errorf("there were unfulfilled expectations: %s", err)
	// 		return
	// 	}
	// 	// if !reflect.DeepEqual(user, u) {
	// 	// 	t.Errorf("results not match, want %v, have %v", u, user)
	// 	// 	return
	// 	// }
	// })
	t.Run("error insert tags", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description", "imgs"}).
			AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", "всем привет", pq.Array([]string{"img1", "img2"}))
		mock.ExpectQuery("update").WithArgs(
			"Ilyagu",
			"valid@valid.ru",
			"2001-06-29",
			"всем привет",
			pq.Array([]string{"img1", "img2"}),
		).WillReturnRows(rows)

		vals := []driver.Value{1, "anime", "music"}

		// mock.ExpectQuery("delete").WithArgs(1).WillReturnError(sql.ErrNoRows)
		mock.ExpectQuery("insert into profile_tag").WithArgs(vals).WillReturnError(sql.ErrNoRows)

		// rowsTags := sqlmock.NewRows([]string{"tag_name"}).AddRow("tag1").AddRow("tag2")
		// mock.ExpectQuery("select tag_name").WithArgs(1).WillReturnRows(rowsTags)

		_, err := repo.UpdateUser(context.TODO(), u)

		if err == nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		// if !reflect.DeepEqual(user, u) {
		// 	t.Errorf("results not match, want %v, have %v", u, user)
		// 	return
		// }
	})
}
