package Handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"server/MockDB"
	"server/Models"
	"testing"
)

const (
	correctCase = iota + 1
	wrongCase
)

type TestCase struct {
	testType   int
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
				Name:  "sessionId",
				Value: "123",
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":{"id":1,"email":"testCurrentUser1@mail.ru"}}`,
		},
		TestCase{
			BodyReq: nil,
			CookieReq: http.Cookie{
				Name:  "sessionId",
				Value: "case wrong cookie",
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
		TestCase{
			BodyReq:    nil,
			CookieReq:  http.Cookie{},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
	}

	testDB := MockDB.NewMockDB()
	testSessionDB := MockDB.NewSessionDB()

	env := &Env{
		DB:        testDB,
		SessionDB: testSessionDB,
	}
	testDB.CreateUser(Models.LoginUser{
		Email:    "testCurrentUser1@mail.ru",
		Password: "123456qQ",
	})
	testSessionDB.NewSessionCookie("123", 1)

	for caseNum, item := range cases {
		r := httptest.NewRequest("GET", "/api/v1/currentuser", item.BodyReq)
		r.AddCookie(&item.CookieReq)
		w := httptest.NewRecorder()

		env.CurrentUser(w, r)

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
			testType:   correctCase,
			BodyReq:    bytes.NewReader([]byte(`{"email":"testLogin1@mail.ru","password":"123456qQ"}`)),
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":{"id":1,"email":"testLogin1@mail.ru"}}`,
		},
		TestCase{
			testType:   wrongCase,
			BodyReq:    bytes.NewReader([]byte(`wrong input data`)),
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":400,"body":null}`,
		},
		TestCase{
			testType:   wrongCase,
			BodyReq:    bytes.NewReader([]byte(`{"email":"wrongEmail","password":"wrongPassword"}`)),
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
		TestCase{
			testType:   wrongCase,
			BodyReq:    bytes.NewReader([]byte(`{"email":"testLogin1@mail.ru","password":"wrongPassword"}`)),
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
	}

	testDB := MockDB.NewMockDB()
	testSessionDB := MockDB.NewSessionDB()

	env := &Env{
		DB:        testDB,
		SessionDB: testSessionDB,
	}
	testDB.CreateUser(Models.LoginUser{
		Email:    "testLogin1@mail.ru",
		Password: "123456qQ",
	})

	for caseNum, item := range cases {
		r := httptest.NewRequest("POST", "/api/v1/login", item.BodyReq)
		r.AddCookie(&item.CookieReq)
		w := httptest.NewRecorder()

		env.LoginHandler(w, r)

		if w.Code != item.StatusCode {
			t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
				caseNum+1, w.Code, item.StatusCode)
		}

		if w.Body.String() != item.BodyResp {
			t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
				caseNum+1, w.Body.String(), item.BodyResp)
		}

		if item.testType == correctCase {
			if !testSessionDB.IsSessionByUserID(1) {
				t.Errorf("TestCase [%d]:\nsession was not created", caseNum+1)
			}
			testSessionDB.DropCookies()
		}
	}
}

func TestSignup(t *testing.T) {
	t.Parallel()

	email := "testSignup1@mail.ru"
	password := "123456qQ"

	expectedUser := Models.MakeUser(2, email, password)

	cases := []TestCase{
		TestCase{
			testType:   correctCase,
			BodyReq:    bytes.NewReader([]byte(`{"email":"` + email + `","password":"` + password + `"}`)),
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":null}`,
		},
		TestCase{
			testType:   wrongCase,
			BodyReq:    bytes.NewReader([]byte(`wrong input data`)),
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":400,"body":null}`,
		},
		TestCase{
			testType:   wrongCase,
			BodyReq:    bytes.NewReader([]byte(`{"email":"firsUser@mail.ru","password":"EmailAlreadyExists"}`)),
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":1001,"body":null}`,
		},
	}

	testDB := MockDB.NewMockDB()
	testSessionDB := MockDB.NewSessionDB()

	env := &Env{
		DB:        testDB,
		SessionDB: testSessionDB,
	}

	for caseNum, item := range cases {
		testDB.DropUsers()
		testDB.CreateUser(Models.LoginUser{
			Email:    "firsUser@mail.ru",
			Password: "123456qQ",
		})

		r := httptest.NewRequest("POST", "/api/v1/signup", item.BodyReq)
		r.AddCookie(&item.CookieReq)
		w := httptest.NewRecorder()

		env.SignupHandler(w, r)

		if w.Code != item.StatusCode {
			t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
				caseNum+1, w.Code, item.StatusCode)
		}

		if w.Body.String() != item.BodyResp {
			t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
				caseNum+1, w.Body.String(), item.BodyResp)
		}

		if item.testType == correctCase {
			if !testSessionDB.IsSessionByUserID(expectedUser.ID) {
				t.Errorf("TestCase [%d]:\nsession was not created", caseNum+1)
			}
			testSessionDB.DropCookies()

			newUser, _ := testDB.GetUser(email)
			if !reflect.DeepEqual(newUser, expectedUser) {
				t.Errorf("TestCase [%d]:\nuser was not created", caseNum+1)
			}
		}
	}
}

func TestLogout(t *testing.T) {
	t.Parallel()
	cases := []TestCase{
		TestCase{
			testType: correctCase,
			CookieReq: http.Cookie{
				Name:  "sessionId",
				Value: "123",
			},
			StatusCode: http.StatusOK,
		},
		TestCase{
			testType: wrongCase,
			CookieReq: http.Cookie{
				Name:  "sessionId",
				Value: "case wrong cookie",
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":500,"body":null}`,
		},
		TestCase{
			testType:   wrongCase,
			CookieReq:  http.Cookie{},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
	}

	testDB := MockDB.NewMockDB()
	testSessionDB := MockDB.NewSessionDB()

	env := &Env{
		DB:        testDB,
		SessionDB: testSessionDB,
	}
	user, _ := testDB.CreateUser(Models.LoginUser{
		Email:    "testLogout1@mail.ru",
		Password: "123456qQ",
	})

	for caseNum, item := range cases {
		r := httptest.NewRequest("GET", "/api/v1/logout", item.BodyReq)
		r.AddCookie(&item.CookieReq)
		w := httptest.NewRecorder()

		testSessionDB.NewSessionCookie("123", user.ID)

		env.LogoutHandler(w, r)

		if w.Code != item.StatusCode {
			t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
				caseNum+1, w.Code, item.StatusCode)
		}

		if w.Body.String() != item.BodyResp {
			t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
				caseNum+1, w.Body.String(), item.BodyResp)
		}

		if testSessionDB.IsSessionByUserID(user.ID) && item.testType == correctCase {
			t.Errorf("TestCase [%d]:\nuser session not ended", caseNum+1)
		}
	}
}

func TestNextUser(t *testing.T) {
	t.Parallel()
	cases := []TestCase{
		TestCase{
			testType: correctCase,
			BodyReq:  bytes.NewReader([]byte(`{"id":321}`)),
			CookieReq: http.Cookie{
				Name:  "sessionId",
				Value: "123",
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":200,"body":{"id":1,"email":"testNextUser1@mail.ru"}}`,
		},
		TestCase{
			testType: wrongCase,
			BodyReq:  bytes.NewReader([]byte(`{"id":321}`)),
			CookieReq: http.Cookie{
				Name:  "sessionId",
				Value: "case wrong cookie",
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
		TestCase{
			testType:   wrongCase,
			BodyReq:    bytes.NewReader([]byte(`{"id":321}`)),
			CookieReq:  http.Cookie{},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
		TestCase{
			testType: wrongCase,
			BodyReq:  bytes.NewReader([]byte("wrong json")),
			CookieReq: http.Cookie{
				Name:  "sessionId",
				Value: "123",
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":400,"body":null}`,
		},
		TestCase{
			testType: wrongCase,
			BodyReq:  bytes.NewReader([]byte(`{"id":1}`)),
			CookieReq: http.Cookie{
				Name:  "sessionId",
				Value: "123",
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
	}

	testDB := MockDB.NewMockDB()
	testSessionDB := MockDB.NewSessionDB()

	env := &Env{
		DB:        testDB,
		SessionDB: testSessionDB,
	}

	testDB.CreateUser(Models.LoginUser{
		Email:    "testNextUser1@mail.ru",
		Password: "123456qQ\"",
	})

	currenUser, _ := testDB.CreateUser(Models.LoginUser{
		Email:    "testCurrUser1@mail.ru",
		Password: "123456qQ\"",
	})
	testSessionDB.NewSessionCookie("123", currenUser.ID)

	for caseNum, item := range cases {
		r := httptest.NewRequest("POST", "/api/v1/nextswipeuser", item.BodyReq)
		r.AddCookie(&item.CookieReq)
		w := httptest.NewRecorder()

		env.NextUserHandler(w, r)

		if w.Code != item.StatusCode {
			t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
				caseNum+1, w.Code, item.StatusCode)
		}

		if w.Body.String() != item.BodyResp {
			t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
				caseNum+1, w.Body.String(), item.BodyResp)
		}

		if !testDB.IsSwiped(currenUser.ID, 321) && item.testType == correctCase {
			t.Errorf("TestCase [%d]:\nswipe not saved", caseNum+1)
		}
		testDB.DropSwipes()
	}
}

func TestEditProfile(t *testing.T) {
	t.Parallel()

	requestUser := Models.User{
		Name:        "testEdit",
		Date:        "1999-10-25",
		Description: "Description Description Description Description",
		ImgSrc:      "/img/testEdit/",
		Tags:        []string{"Tags", "Tags", "Tags", "Tags", "Tags"},
	}
	bodyReq, _ := json.Marshal(requestUser)

	expectedUser := Models.MakeUser(1, "testEdit@mail.ru", "123456qQ")
	expectedUser.FillProfile(Models.User{
		Name:        requestUser.Name,
		Date:        requestUser.Date,
		Description: requestUser.Description,
		ImgSrc:      requestUser.ImgSrc,
		Tags:        requestUser.Tags,
	})

	BodyRespByte, _ := json.Marshal(Models.JSON{
		Status: StatusOK,
		Body:   expectedUser,
	})

	cases := []TestCase{
		TestCase{
			testType: correctCase,
			BodyReq:  bytes.NewReader(bodyReq),
			CookieReq: http.Cookie{
				Name:  "sessionId",
				Value: "123",
			},
			StatusCode: http.StatusOK,
			BodyResp:   string(BodyRespByte),
		},
		TestCase{
			testType: wrongCase,
			BodyReq:  bytes.NewReader(bodyReq),
			CookieReq: http.Cookie{
				Name:  "sessionId",
				Value: "case wrong cookie",
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
		TestCase{
			testType:   wrongCase,
			BodyReq:    bytes.NewReader(bodyReq),
			CookieReq:  http.Cookie{},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":404,"body":null}`,
		},
		TestCase{
			testType: wrongCase,
			BodyReq:  bytes.NewReader([]byte(`{"name":"testEdit","date":"wrong-format-data","description":"Description Description Description Description","imgSrc":"/img/testEdit/","tags":["Tags","Tags","Tags","Tags","Tags"]}`)),
			CookieReq: http.Cookie{
				Name:  "sessionId",
				Value: "123",
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":400,"body":null}`,
		},
		TestCase{
			testType: wrongCase,
			BodyReq:  bytes.NewReader([]byte(`wrong data`)),
			CookieReq: http.Cookie{
				Name:  "sessionId",
				Value: "123",
			},
			StatusCode: http.StatusOK,
			BodyResp:   `{"status":400,"body":null}`,
		},
	}

	testDB := MockDB.NewMockDB()
	testSessionDB := MockDB.NewSessionDB()

	env := &Env{
		DB:        testDB,
		SessionDB: testSessionDB,
	}

	for caseNum, item := range cases {
		testDB.DropUsers()
		testSessionDB.DropCookies()
		currenUser, _ := testDB.CreateUser(Models.LoginUser{
			Email:    expectedUser.Email,
			Password: "123456qQ",
		})
		testSessionDB.NewSessionCookie("123", currenUser.ID)

		r := httptest.NewRequest("POST", "/api/v1/editprofile", item.BodyReq)
		r.AddCookie(&item.CookieReq)
		w := httptest.NewRecorder()

		env.EditProfileHandler(w, r)

		if w.Code != item.StatusCode {
			t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
				caseNum+1, w.Code, item.StatusCode)
		}

		if w.Body.String() != item.BodyResp {
			t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
				caseNum+1, w.Body.String(), item.BodyResp)
		}

		if item.testType == correctCase {
			updateUser, err := testDB.GetUser(currenUser.Email)
			if err != nil {
				t.Errorf("TestCase [%d]:\nprofile was not created", caseNum+1)
			}
			if !reflect.DeepEqual(updateUser, expectedUser) {
				t.Errorf("TestCase [%d]:\nwrong profile: \ngot %v\nexpected %v",
					caseNum+1, updateUser, expectedUser)
			}
		}
	}
}