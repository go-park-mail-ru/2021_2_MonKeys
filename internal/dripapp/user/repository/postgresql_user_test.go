package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"dripapp/configs"
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

		rowsTags := sqlmock.NewRows([]string{"tagname"}).
			AddRow("tag1").
			AddRow("tag2")
		mock.ExpectQuery("select tagname").WithArgs(1).WillReturnRows(rowsTags)

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

		mock.ExpectQuery("select tagname").WithArgs(1).WillReturnError(sql.ErrNoRows)

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

		mock.ExpectQuery("select tagname").WithArgs(1).WillReturnError(fmt.Errorf("db_error"))

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

		rowsTags := sqlmock.NewRows([]string{"tagname"}).
			AddRow("tag1").
			AddRow("tag2")
		mock.ExpectQuery("select tagname").WithArgs(1).WillReturnRows(rowsTags)

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

		mock.ExpectQuery("select tagname").WithArgs(1).WillReturnError(sql.ErrNoRows)

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

		mock.ExpectQuery("select tagname").WithArgs(1).WillReturnError(fmt.Errorf("db_error"))

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

func TestInsertTags(t *testing.T) {
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

	vals := []driver.Value{1, "anime", "music"}

	t.Run("good insert tags", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("insert").WithArgs(vals...).WillReturnRows(rows)

		err := repo.insertTags(context.TODO(), 1, []string{"anime", "music"})

		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	})
	t.Run("len tags nil", func(t *testing.T) {
		err := repo.insertTags(context.TODO(), 1, []string{})

		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
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
		Tags:        []string{"anime", "music"},
	}

	vals := []driver.Value{1, "anime", "music"}

	t.Run("good update", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description", "imgs"}).
			AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", "всем привет", pq.Array([]string{"img1", "img2"}))
		mock.ExpectQuery("update").WithArgs(
			"Ilyagu",
			"valid@valid.ru",
			"2001-06-29",
			"всем привет",
			pq.Array([]string{"img1", "img2"}),
		).WillReturnRows(rows)

		rowsDel := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("delete").WithArgs(1).WillReturnRows(rowsDel)

		rowsInsTags := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("insert").WithArgs(vals...).WillReturnRows(rowsInsTags)

		rowsTags := sqlmock.NewRows([]string{"tagname"}).AddRow("anime").AddRow("music")
		mock.ExpectQuery("select tagname").WithArgs(1).WillReturnRows(rowsTags)

		user, err := repo.UpdateUser(context.TODO(), u)

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
	t.Run("error update", func(t *testing.T) {
		mock.ExpectQuery("update").WithArgs(
			"Ilyagu",
			"valid@valid.ru",
			"2001-06-29",
			"всем привет",
			pq.Array([]string{"img1", "img2"}),
		).WillReturnError(sql.ErrNoRows)

		_, err = repo.UpdateUser(context.TODO(), u)

		if err == nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	})
	t.Run("error delete tags", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description", "imgs"}).
			AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", "всем привет", pq.Array([]string{"img1", "img2"}))
		mock.ExpectQuery("update").WithArgs(
			"Ilyagu",
			"valid@valid.ru",
			"2001-06-29",
			"всем привет",
			pq.Array([]string{"img1", "img2"}),
		).WillReturnRows(rows)

		mock.ExpectQuery("delete").WithArgs(1).WillReturnError(sql.ErrNoRows)

		_, err := repo.UpdateUser(context.TODO(), u)

		if err == nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	})
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

		rowsDel := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("delete").WithArgs(1).WillReturnRows(rowsDel)

		mock.ExpectQuery("insert").WithArgs(vals...).WillReturnError(sql.ErrNoRows)

		_, err = repo.UpdateUser(context.TODO(), u)

		if err == nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	})
	t.Run("error get tags", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description", "imgs"}).
			AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", "всем привет", pq.Array([]string{"img1", "img2"}))
		mock.ExpectQuery("update").WithArgs(
			"Ilyagu",
			"valid@valid.ru",
			"2001-06-29",
			"всем привет",
			pq.Array([]string{"img1", "img2"}),
		).WillReturnRows(rows)

		rowsDel := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("delete").WithArgs(1).WillReturnRows(rowsDel)

		rowsInsTags := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("insert").WithArgs(vals...).WillReturnRows(rowsInsTags)

		mock.ExpectQuery("select tagname").WithArgs(1).WillReturnError(fmt.Errorf("some error"))

		_, err = repo.UpdateUser(context.TODO(), u)

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

func TestGetTags(t *testing.T) {
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

	tags := make(map[uint64]string)
	tags[0] = "anime"
	tags[1] = "music"

	t.Run("good get tags", func(t *testing.T) {
		rowsTags := sqlmock.NewRows([]string{"tagname"}).
			AddRow("anime").
			AddRow("music")
		mock.ExpectQuery("select tagname").WillReturnRows(rowsTags)

		testTags, err := repo.GetTags(context.TODO())

		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if !reflect.DeepEqual(testTags, tags) {
			t.Errorf("results not match, want \n%v\n, have \n%v\n", testTags, tags)
			return
		}
	})
	t.Run("error get tags", func(t *testing.T) {
		mock.ExpectQuery("select tagname").WillReturnError(sql.ErrNoRows)

		_, err := repo.GetTags(context.TODO())

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

func TestUpdateImgs(t *testing.T) {
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

	t.Run("good update imgs", func(t *testing.T) {
		rowsTags := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("update profile").WithArgs(1, pq.Array([]string{"img1", "img2"})).WillReturnRows(rowsTags)

		err := repo.UpdateImgs(context.TODO(), 1, []string{"img1", "img2"})

		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	})
	t.Run("error update imgs", func(t *testing.T) {
		mock.ExpectQuery("update profile").WithArgs(1, pq.Array([]string{"img1", "img2"})).WillReturnError(sql.ErrNoRows)

		err := repo.UpdateImgs(context.TODO(), 1, []string{"img1", "img2"})

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

func TestAddReaction(t *testing.T) {
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

	t.Run("good add reaction", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("insert into reactions").WithArgs(1, 1, 1).WillReturnRows(rows)

		err := repo.AddReaction(context.TODO(), 1, 1, 1)

		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	})
	t.Run("error add reaction", func(t *testing.T) {
		mock.ExpectQuery("insert into reactions").WithArgs(1, 1, 1).WillReturnError(sql.ErrNoRows)

		err := repo.AddReaction(context.TODO(), 1, 1, 1)

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

func TestGetNext(t *testing.T) {
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

	us := []models.User{
		{
			ID:          1,
			Email:       "valid@valid.ru",
			Password:    "!Nagdimaev2001",
			Name:        "Ilyagu",
			Gender:      "male",
			Prefer:      "female",
			Date:        "2001-06-29",
			Description: " я хач",
			Age:         "20",
			Imgs:        []string{"img1", "img2"},
			Tags:        []string{"tag1", "tag2"},
		},
	}

	t.Run("good get next", func(t *testing.T) {
		row := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description"}).
			AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", " я хач")
		mock.ExpectQuery("select op.id").WithArgs(1).WillReturnRows(row)

		rowImgs := sqlmock.NewRows([]string{"imgs"}).AddRow(pq.Array([]string{"img1", "img2"}))
		mock.ExpectQuery("SELECT imgs").WithArgs(1).WillReturnRows(rowImgs)

		rowsTags := sqlmock.NewRows([]string{"tagname"}).
			AddRow("tag1").
			AddRow("tag2")
		mock.ExpectQuery("select tagname").WithArgs(1).WillReturnRows(rowsTags)

		users, err := repo.GetNextUserForSwipe(context.TODO(), 1, "")

		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if !reflect.DeepEqual(users, us) {
			t.Errorf("results not match, want %v, have %v", users, us)
			return
		}
	})
	t.Run("error get next", func(t *testing.T) {
		mock.ExpectQuery("select op.id").WithArgs(1).WillReturnError(sql.ErrNoRows)

		_, err := repo.GetNextUserForSwipe(context.TODO(), 1, "")

		if err == nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	})
	t.Run("error get imgs", func(t *testing.T) {
		row := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description"}).
			AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", " я хач")
		mock.ExpectQuery("select op.id").WithArgs(1).WillReturnRows(row)

		mock.ExpectQuery("SELECT imgs").WithArgs(1).WillReturnError(fmt.Errorf("some error"))

		_, err := repo.GetNextUserForSwipe(context.TODO(), 1, "")

		if err == nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	})
	t.Run("error get tags", func(t *testing.T) {
		row := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description"}).
			AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", " я хач")
		mock.ExpectQuery("select op.id").WithArgs(1).WillReturnRows(row)

		rowImgs := sqlmock.NewRows([]string{"imgs"}).AddRow(pq.Array([]string{"img1", "img2"}))
		mock.ExpectQuery("SELECT imgs").WithArgs(1).WillReturnRows(rowImgs)

		mock.ExpectQuery("select tagname").WithArgs(1).WillReturnError(sql.ErrNoRows)

		_, err := repo.GetNextUserForSwipe(context.TODO(), 1, "")

		if err == nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	})
	t.Run("error date", func(t *testing.T) {
		row := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description"}).
			AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "fsdgsdg", " я хач")
		mock.ExpectQuery("select op.id").WithArgs(1).WillReturnRows(row)

		_, err := repo.GetNextUserForSwipe(context.TODO(), 1, "")

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

func TestGetNextMatches(t *testing.T) {
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

	us := []models.User{
		{
			ID:          1,
			Name:        "Ilyagu",
			Email:       "valid@valid.ru",
			Password:    "!Nagdimaev2001",
			Date:        "2001-06-29",
			Description: " я хач",
			Age:         "20",
			Imgs:        []string{"img1", "img2"},
			Tags:        []string{"tag1", "tag2"},
		},
	}

	t.Run("good get next matches", func(t *testing.T) {
		row := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description"}).
			AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", " я хач")
		mock.ExpectQuery("select op.id").WithArgs(1).WillReturnRows(row)

		rowImgs := sqlmock.NewRows([]string{"imgs"}).AddRow(pq.Array([]string{"img1", "img2"}))
		mock.ExpectQuery("SELECT imgs").WithArgs(1).WillReturnRows(rowImgs)

		rowsTags := sqlmock.NewRows([]string{"tagname"}).
			AddRow("tag1").
			AddRow("tag2")
		mock.ExpectQuery("select tagname").WithArgs(1).WillReturnRows(rowsTags)

		users, err := repo.GetUsersMatches(context.TODO(), 1)

		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if !reflect.DeepEqual(users, us) {
			t.Errorf("results not match, want %v, have %v", users, us)
			return
		}
	})
	t.Run("error get next", func(t *testing.T) {
		mock.ExpectQuery("select op.id").WithArgs(1).WillReturnError(sql.ErrNoRows)

		_, err := repo.GetUsersMatches(context.TODO(), 1)

		if err == nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	})
	t.Run("error get imgs", func(t *testing.T) {
		row := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description"}).
			AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", " я хач")
		mock.ExpectQuery("select op.id").WithArgs(1).WillReturnRows(row)

		mock.ExpectQuery("SELECT imgs").WithArgs(1).WillReturnError(fmt.Errorf("some error"))

		_, err := repo.GetUsersMatches(context.TODO(), 1)

		if err == nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	})
	t.Run("error get tags", func(t *testing.T) {
		row := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description"}).
			AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "2001-06-29", " я хач")
		mock.ExpectQuery("select op.id").WithArgs(1).WillReturnRows(row)

		rowImgs := sqlmock.NewRows([]string{"imgs"}).AddRow(pq.Array([]string{"img1", "img2"}))
		mock.ExpectQuery("SELECT imgs").WithArgs(1).WillReturnRows(rowImgs)

		mock.ExpectQuery("select tagname").WithArgs(1).WillReturnError(sql.ErrNoRows)

		_, err := repo.GetUsersMatches(context.TODO(), 1)

		if err == nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	})
	t.Run("error date", func(t *testing.T) {
		row := sqlmock.NewRows([]string{"id", "name", "email", "password", "date", "description"}).
			AddRow(1, "Ilyagu", "valid@valid.ru", "!Nagdimaev2001", "fsdgsdg", " я хач")
		mock.ExpectQuery("select op.id").WithArgs(1).WillReturnRows(row)

		_, err := repo.GetUsersMatches(context.TODO(), 1)

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

func TestGetLikes(t *testing.T) {
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

	testLikes := []uint64{1, 2, 10}

	t.Run("good get likes", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id1"}).AddRow(1).AddRow(2).AddRow(10)
		mock.ExpectQuery("select").WithArgs(1).WillReturnRows(rows)

		likes, err := repo.GetLikes(context.TODO(), 1)

		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if !reflect.DeepEqual(testLikes, likes) {
			t.Errorf("results not match, want %v, have %v", testLikes, likes)
			return
		}
	})
	t.Run("error get likes", func(t *testing.T) {
		mock.ExpectQuery("select").WithArgs(1).WillReturnError(sql.ErrNoRows)

		_, err := repo.GetLikes(context.TODO(), 1)

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

func TestDeleteLikes(t *testing.T) {
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

	t.Run("good delete likes", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("delete").WithArgs(1, 1).WillReturnRows(rows)

		err := repo.DeleteLike(context.TODO(), 1, 1)

		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	})
	t.Run("error delete likes", func(t *testing.T) {
		mock.ExpectQuery("delete").WithArgs(1, 1).WillReturnError(sql.ErrNoRows)

		err := repo.DeleteLike(context.TODO(), 1, 1)

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

func TestAddMatch(t *testing.T) {
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

	t.Run("good add match", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("insert into matches").WithArgs(1, 1).WillReturnRows(rows)

		err := repo.AddMatch(context.TODO(), 1, 1)

		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	})
	t.Run("error add match", func(t *testing.T) {
		mock.ExpectQuery("insert into matches").WithArgs(1, 1).WillReturnError(sql.ErrNoRows)

		err := repo.AddMatch(context.TODO(), 1, 1)

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

func TestNew(t *testing.T) {
	t.Parallel()
	t.Run("success new", func(t *testing.T) {
		_, err := NewPostgresUserRepository(configs.Postgres)
		if err != nil {
			t.Error()
		}
	})
	t.Run("error new", func(t *testing.T) {
		_, err := NewPostgresUserRepository(configs.PostgresConfig{
			User:     "flksdmflksdklf",
			Password: "fsdmflsldfmlsdf",
			DBName:   "f;lsd,fls,df;",
		})
		if err != nil {
			t.Error()
		}
	})
}
