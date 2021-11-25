package repository

import (
	"context"
	"database/sql"
	"dripapp/configs"
	"dripapp/internal/microservices/chat/models"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestChats(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := &PostgreChatRepo{
		Conn: *sqlxDB,
	}

	testTime := time.Now()

	testMessage := []models.Message{
		{
			MessageID: 1,
			FromID:    1,
			ToID:      2,
			Text:      "text",
			Date:      testTime,
		},
	}

	testChats := []models.Chat{
		{
			FromUserID: 1,
			Name:       "Ilyagu",
			Img:        "lol",
			Messages:   testMessage,
		},
	}

	t.Run("good get chats", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "img"}).AddRow("1", "Ilyagu", "lol")
		mock.ExpectQuery("select").WithArgs(1).WillReturnRows(rows)

		messageRows := sqlmock.NewRows([]string{"message_id", "from_id", "to_id", "text", "date"}).
			AddRow(1, 1, 2, "text", testTime)
		mock.ExpectQuery("select").WithArgs(1, 1).WillReturnRows(messageRows)

		chats, err := repo.GetChats(context.TODO(), 1)

		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if !reflect.DeepEqual(chats, testChats) {
			t.Errorf("results not match, want %v, have %v", chats, testChats)
			return
		}
	})
	t.Run("error get chats", func(t *testing.T) {
		mock.ExpectQuery("select").WithArgs(1).WillReturnError(sql.ErrNoRows)

		_, err := repo.GetChats(context.TODO(), 1)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if err == nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
	})
	t.Run("error get last message", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "img"}).AddRow("1", "Ilyagu", "lol")
		mock.ExpectQuery("select").WithArgs(1).WillReturnRows(rows)

		mock.ExpectQuery("select").WithArgs(1, 1).WillReturnError(sql.ErrNoRows)

		_, err := repo.GetChats(context.TODO(), 1)

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

func TestChat(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := &PostgreChatRepo{
		Conn: *sqlxDB,
	}

	testTime := time.Now()

	testMessage := []models.Message{
		{
			MessageID: 1,
			FromID:    1,
			ToID:      2,
			Text:      "text",
			Date:      testTime,
		},
		{
			MessageID: 1,
			FromID:    1,
			ToID:      2,
			Text:      "text",
			Date:      testTime,
		},
	}

	t.Run("good get chat", func(t *testing.T) {
		messageRows := sqlmock.NewRows([]string{"message_id", "from_id", "to_id", "text", "date"}).
			AddRow(1, 1, 2, "text", testTime).
			AddRow(1, 1, 2, "text", testTime)
		mock.ExpectQuery("select").WithArgs(1, 1, 1).WillReturnRows(messageRows)

		messages, err := repo.GetChat(context.TODO(), 1, 1, 1)

		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if !reflect.DeepEqual(messages, testMessage) {
			t.Errorf("results not match, want %v, have %v", messages, testMessage)
			return
		}
	})
	t.Run("error get chat", func(t *testing.T) {
		mock.ExpectQuery("select").WithArgs(1, 1, 1).WillReturnError(sql.ErrNoRows)

		_, err := repo.GetChat(context.TODO(), 1, 1, 1)

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

func TestSaveMessage(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := &PostgreChatRepo{
		Conn: *sqlxDB,
	}

	testTime := time.Now()

	testMessage := models.Message{
		MessageID: 1,
		FromID:    1,
		ToID:      2,
		Text:      "text",
		Date:      testTime,
	}

	t.Run("good save message", func(t *testing.T) {
		messageRows := sqlmock.NewRows([]string{"message_id", "from_id", "to_id", "text", "date"}).
			AddRow(1, 1, 2, "text", testTime)
		mock.ExpectQuery("insert").WithArgs(1, 1, "text").WillReturnRows(messageRows)

		message, err := repo.SaveMessage(1, 1, "text")

		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if !reflect.DeepEqual(message, testMessage) {
			t.Errorf("results not match, want %v, have %v", message, testMessage)
			return
		}
	})
	t.Run("error save message", func(t *testing.T) {
		mock.ExpectQuery("insert").WithArgs(1, 1, "text").WillReturnError(sql.ErrNoRows)

		_, err := repo.SaveMessage(1, 1, "text")

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

func TestNew(t *testing.T) {
	t.Parallel()
	t.Run("success new", func(t *testing.T) {
		_, err := NewPostgresChatRepository(configs.Postgres)
		if err != nil {
			t.Error()
		}
	})
	t.Run("error new", func(t *testing.T) {
		_, err := NewPostgresChatRepository(configs.PostgresConfig{
			User:     "flksdmflksdklf",
			Password: "fsdmflsldfmlsdf",
			DBName:   "f;lsd,fls,df;",
		})
		if err != nil {
			t.Error()
		}
	})
}
