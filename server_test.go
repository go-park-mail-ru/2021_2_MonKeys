package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const (
	correctCase = iota + 1
	wrongCase
)

type TestCase struct {
	testType    int
	BodyReq    io.Reader
	CookieReq  http.Cookie
	StatusCode int
	BodyResp   string
}

func TestCurrentUser(t *testing.T) {
	t.Parallel()
	cases := []TestCase{
		TestCase{
			BodyReq: nil,
			CookieReq: http.Cookie{
				Name:     "sessionId",
				Value:    "123",
			},
			StatusCode: http.StatusOK,
			BodyResp: `{"status":200,"body":{"id":1,"name":"","email":"testCurrentUser1@mail.ru","age":0,"description":"","imgSrc":"","tags":null}}`,
		},
		TestCase{
			BodyReq: nil,
			CookieReq: http.Cookie{
				Name:     "sessionId",
				Value:    "123123",
			},
			StatusCode: http.StatusOK,
			BodyResp: `{"status":404,"body":null}`,
		},
		TestCase{
			BodyReq: nil,
			CookieReq: http.Cookie{},
			StatusCode: http.StatusOK,
			BodyResp: `{"status":404,"body":null}`,
		},
	}

	testDB := NewMockDB()
	testSessionDB := NewSessionDB()

	env := &Env{
		db:        testDB,
		sessionDB: testSessionDB,
	}
	testDB.createUser(LoginUser{
		Email: "testCurrentUser1@mail.ru",
		Password: "123456qQ",
	})
	testSessionDB.cookies["123"] = 1

	for caseNum, item := range cases {
		r := httptest.NewRequest("GET", "/api/v1/currentuser", item.BodyReq)
		r.AddCookie(&item.CookieReq)
		w := httptest.NewRecorder()

		env.currentUser(w, r)

		if w.Code != item.StatusCode {
			t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
				caseNum+1, w.Code, item.StatusCode)
		}

		if w.Body.String() != item.BodyResp {
			t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
				caseNum+1, w.Body.String(), item.BodyResp)
		}
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()
	cases := []TestCase{
		TestCase{
			testType: correctCase,
			BodyReq: bytes.NewReader([]byte(`{"email":"testLogin1@mail.ru","password":"123456qQ"}`)),
			StatusCode: http.StatusOK,
			BodyResp: `{"status":200,"body":null}`,
		},
		TestCase{
			testType: wrongCase,
			BodyReq: bytes.NewReader([]byte(`wrong input data`)),
			StatusCode: http.StatusOK,
			BodyResp: `{"status":400,"body":null}`,
		},
		TestCase{
			testType: wrongCase,
			BodyReq: bytes.NewReader([]byte(`{"email":"wrongEmail","password":"wrongPassword"}`)),
			StatusCode: http.StatusOK,
			BodyResp: `{"status":404,"body":null}`,
		},
		TestCase{
			testType: wrongCase,
			BodyReq: bytes.NewReader([]byte(`{"email":"testLogin1@mail.ru","password":"wrongPassword"}`)),
			StatusCode: http.StatusOK,
			BodyResp: `{"status":404,"body":null}`,
		},
	}

	testDB := NewMockDB()
	testSessionDB := NewSessionDB()

	env := &Env{
		db:        testDB,
		sessionDB: testSessionDB,
	}
	testDB.createUser(LoginUser{
		Email: "testLogin1@mail.ru",
		Password: "123456qQ",
	})

	for caseNum, item := range cases {
		r := httptest.NewRequest("POST", "/api/v1/login", item.BodyReq)
		r.AddCookie(&item.CookieReq)
		w := httptest.NewRecorder()

		env.loginHandler(w, r)

		if w.Code != item.StatusCode {
			t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
				caseNum+1, w.Code, item.StatusCode)
		}

		if !testSessionDB.isSessionByUserID(1) && item.testType == correctCase {
			t.Errorf("TestCase [%d]:\nsession was not created", caseNum+1)
		}
		testSessionDB.cookies = make(map[string]uint64)

		if w.Body.String() != item.BodyResp {
			t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
				caseNum+1, w.Body.String(), item.BodyResp)
		}
	}
}

func TestSignup(t *testing.T) {
	t.Parallel()

	email := "testSignup1@mail.ru"
	password := "123456qQ"

	newUserID := uint64(2)
	expectedID := newUserID
	expectedUsers := makeUser(expectedID, email, password)

	cases := []TestCase{
		TestCase{
			testType: correctCase,
			BodyReq: bytes.NewReader([]byte(`{"email":"` + email + `","password":"` + password + `"}`)),
			StatusCode: http.StatusOK,
			BodyResp: `{"status":200,"body":null}`,
		},
		TestCase{
			testType: wrongCase,
			BodyReq: bytes.NewReader([]byte(`wrong input data`)),
			StatusCode: http.StatusOK,
			BodyResp: `{"status":400,"body":null}`,
		},
		TestCase{
			testType: wrongCase,
			BodyReq: bytes.NewReader([]byte(`{"email":"firsUser@mail.ru","password":"EmailAlreadyExists"}`)),
			StatusCode: http.StatusOK,
			BodyResp: `{"status":1001,"body":null}`,
		},
	}

	testDB := NewMockDB()
	testSessionDB := NewSessionDB()

	env := &Env{
		db:        testDB,
		sessionDB: testSessionDB,
	}

	for caseNum, item := range cases {
		testDB.users = make(map[uint64]User)
		testDB.createUser(LoginUser{
			Email: "firsUser@mail.ru",
			Password: "123456qQ",
		})

		r := httptest.NewRequest("POST", "/api/v1/signup", item.BodyReq)
		r.AddCookie(&item.CookieReq)
		w := httptest.NewRecorder()

		env.signupHandler(w, r)

		if w.Code != item.StatusCode {
			t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
				caseNum+1, w.Code, item.StatusCode)
		}

		if !testSessionDB.isSessionByUserID(newUserID) && item.testType == correctCase {
			t.Errorf("TestCase [%d]:\nsession was not created", caseNum+1)
		}
		testSessionDB.cookies = make(map[string]uint64)

		newUser, _ := testDB.getUser(email)
		if !reflect.DeepEqual(newUser, expectedUsers) && item.testType == correctCase {
			t.Errorf("TestCase [%d]:\nuser was not created", caseNum+1)
		}

		if w.Body.String() != item.BodyResp {
			t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
				caseNum+1, w.Body.String(), item.BodyResp)
		}
	}
}

func TestLogout(t *testing.T) {
	t.Parallel()
	cases := []TestCase{
		TestCase{
			testType: correctCase,
			CookieReq: http.Cookie{
				Name:     "sessionId",
				Value:    "1234",
			},
			StatusCode: http.StatusOK,
		},
		TestCase{
			testType: wrongCase,
			CookieReq: http.Cookie{
				Name:     "sessionId",
				Value:    "wrongCase cookie",
			},
			StatusCode: http.StatusOK,
			BodyResp: `{"status":500,"body":null}`,
		},
		TestCase{
			testType:   wrongCase,
			CookieReq:  http.Cookie{},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
	}

	testDB := NewMockDB()
	testSessionDB := NewSessionDB()

	env := &Env{
		db:        testDB,
		sessionDB: testSessionDB,
	}
	testDB.createUser(LoginUser{
		Email: "testLogout1@mail.ru",
		Password: "123456qQ",
	})

	for caseNum, item := range cases {
		r := httptest.NewRequest("GET", "/api/v1/logout", item.BodyReq)
		r.AddCookie(&item.CookieReq)
		w := httptest.NewRecorder()

		testSessionDB.cookies["1234"] = 1

		env.logoutHandler(w, r)

		if w.Code != item.StatusCode {
			t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
				caseNum+1, w.Code, item.StatusCode)
		}

		if w.Body.String() != item.BodyResp {
			t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
				caseNum+1, w.Body.String(), item.BodyResp)
		}

		if _, ok := testSessionDB.cookies["1234"]; ok && item.testType == correctCase {
			t.Errorf("TestCase [%d]:\nuser session not ended", caseNum + 1)
		}
	}
}

func TestNextUser(t *testing.T) {
	t.Parallel()
	cases := []TestCase{
		TestCase{
			testType: correctCase,
			BodyReq: bytes.NewReader([]byte(`{"id":321}`)),
			CookieReq: http.Cookie{
				Name:     "sessionId",
				Value:    "123",
			},
			StatusCode: http.StatusOK,
			BodyResp: `{"status":200,"body":{"id":1,"name":"","email":"testNextUser1@mail.ru","age":0,"description":"","imgSrc":"","tags":null}}`,
		},
		TestCase{
			testType: wrongCase,
			BodyReq: bytes.NewReader([]byte(`{"id":321}`)),
			CookieReq: http.Cookie{
				Name:     "sessionId",
				Value:    "123123",
			},
			StatusCode: http.StatusOK,
			BodyResp: `{"status":404,"body":null}`,
		},
		TestCase{
			testType: wrongCase,
			BodyReq: bytes.NewReader([]byte(`{"id":321}`)),
			CookieReq: http.Cookie{},
			StatusCode: http.StatusOK,
			BodyResp: `{"status":404,"body":null}`,
		},
		TestCase{
			testType: wrongCase,
			BodyReq: bytes.NewReader([]byte("wrong json")),
			CookieReq: http.Cookie{
				Name:     "sessionId",
				Value:    "123",
			},
			StatusCode: http.StatusOK,
			BodyResp: `{"status":400,"body":null}`,
		},
		TestCase{
			testType: wrongCase,
			BodyReq: bytes.NewReader([]byte(`{"id":1}`)),
			CookieReq: http.Cookie{
				Name:     "sessionId",
				Value:    "123",
			},
			StatusCode: http.StatusOK,
			BodyResp: `{"status":404,"body":null}`,
		},
	}

	testDB := NewMockDB()
	testSessionDB := NewSessionDB()

	env := &Env{
		db:        testDB,
		sessionDB: testSessionDB,
	}

	testDB.createUser(LoginUser{
		Email: "testNextUser1@mail.ru",
		Password: "123456qQ\"",
	})

	currenUser, _ := testDB.createUser(LoginUser{
		Email: "testCurrUser1@mail.ru",
		Password: "123456qQ\"",
	})
	testSessionDB.cookies["123"] = currenUser.ID

	for caseNum, item := range cases {
		r := httptest.NewRequest("POST", "/api/v1/nextswipeuser", item.BodyReq)
		r.AddCookie(&item.CookieReq)
		w := httptest.NewRecorder()

		env.nextUserHandler(w, r)

		if w.Code != item.StatusCode {
			t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
				caseNum+1, w.Code, item.StatusCode)
		}

		if !testDB.isSwiped(currenUser.ID, 321) && item.testType == correctCase {
			t.Errorf("TestCase [%d]:\nswipe not saved", caseNum+1)
		}
		testDB.swipedUsers = make(map[uint64][]uint64)

		if w.Body.String() != item.BodyResp {
			t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
				caseNum+1, w.Body.String(), item.BodyResp)
		}
	}
}
