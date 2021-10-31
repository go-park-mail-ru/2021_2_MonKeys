package delivery

// import (
// 	"bytes"
// 	"context"
// 	"dripapp/configs"
// 	"dripapp/internal/pkg/models"
// 	"dripapp/internal/pkg/session"
// 	_sessionDelivery "dripapp/internal/pkg/session/delivery"
// 	_userRepo "dripapp/internal/dripapp/user/repository"
// 	_userUCase "dripapp/internal/dripapp/user/usecase"
// 	"encoding/json"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"reflect"
// 	"testing"
// )

// const (
// 	correctCase = iota + 1
// 	wrongCase
// )

// type TestCase struct {
// 	testType   int
// 	BodyReq    io.Reader
// 	CookieReq  http.Cookie
// 	StatusCode int
// 	BodyResp   string
// }

// func TestCurrentUser(t *testing.T) {
// 	t.Parallel()
// 	cases := []TestCase{
// 		{
// 			BodyReq: nil,
// 			CookieReq: http.Cookie{
// 				Name:  "sessionId",
// 				Value: "123",
// 			},
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":200,"body":{"id":1,"email":"testCurrentUser1@mail.ru"}}`,
// 		},
// 		{
// 			BodyReq: nil,
// 			CookieReq: http.Cookie{
// 				Name:  "sessionId",
// 				Value: "case wrong cookie",
// 			},
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":404,"body":null}`,
// 		},
// 		{
// 			BodyReq:    nil,
// 			CookieReq:  http.Cookie{},
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":404,"body":null}`,
// 		},
// 	}

// 	testDB := _userRepo.NewMockDB()
// 	testSessionDB := session.NewSessionDB()

// 	timeoutContext := configs.Timeouts.ContextTimeout
// 	userHandler := &UserHandler{
// 		UserUCase: _userUCase.NewUserUsecase(testDB, testSessionDB, timeoutContext),
// 	}

// 	_, err := testDB.CreateUser(context.TODO(), models.LoginUser{
// 		Email:    "testCurrentUser1@mail.ru",
// 		Password: "123456qQ",
// 	})
// 	if err != nil {
// 		t.Errorf("Create user error")
// 	}
// 	err = testSessionDB.NewSessionCookie("123", 1)
// 	if err != nil {
// 		t.Errorf("New session Cookie error")
// 	}

// 	for caseNum, item := range cases {
// 		r := httptest.NewRequest("GET", "/api/v1/currentuser", item.BodyReq)
// 		r.AddCookie(&item.CookieReq)
// 		w := httptest.NewRecorder()

// 		userHandler.CurrentUser(w, r)

// 		if w.Code != item.StatusCode {
// 			t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
// 				caseNum+1, w.Code, item.StatusCode)
// 		}

// 		if w.Body.String() != item.BodyResp {
// 			t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
// 				caseNum+1, w.Body.String(), item.BodyResp)
// 		}
// 	}
// }

// func TestLogin(t *testing.T) {
// 	t.Parallel()
// 	cases := []TestCase{
// 		{
// 			testType:   correctCase,
// 			BodyReq:    bytes.NewReader([]byte(`{"email":"testLogin1@mail.ru","password":"123456qQ"}`)),
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":200,"body":{"id":1,"email":"testLogin1@mail.ru"}}`,
// 		},
// 		{
// 			testType:   wrongCase,
// 			BodyReq:    bytes.NewReader([]byte(`wrong input data`)),
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":400,"body":null}`,
// 		},
// 		{
// 			testType:   wrongCase,
// 			BodyReq:    bytes.NewReader([]byte(`{"email":"wrongEmail","password":"wrongPassword"}`)),
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":404,"body":null}`,
// 		},
// 		{
// 			testType:   wrongCase,
// 			BodyReq:    bytes.NewReader([]byte(`{"email":"testLogin1@mail.ru","password":"wrongPassword"}`)),
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":404,"body":null}`,
// 		},
// 	}

// 	testDB := _userRepo.NewMockDB()
// 	testSessionDB := session.NewSessionDB()

// 	timeoutContext := configs.Timeouts.ContextTimeout
// 	sessionHandler := &_sessionDelivery.SessionHandler{
// 		UserUCase: _userUCase.NewUserUsecase(testDB, testSessionDB, timeoutContext),
// 	}

// 	_, err := testDB.CreateUser(context.TODO(), models.LoginUser{
// 		Email:    "testLogin1@mail.ru",
// 		Password: "123456qQ",
// 	})
// 	if err != nil {
// 		t.Errorf("Create user failed")
// 	}

// 	for caseNum, item := range cases {
// 		r := httptest.NewRequest("POST", "/api/v1/login", item.BodyReq)
// 		r.AddCookie(&item.CookieReq)
// 		w := httptest.NewRecorder()

// 		sessionHandler.LoginHandler(w, r)

// 		if w.Code != item.StatusCode {
// 			t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
// 				caseNum+1, w.Code, item.StatusCode)
// 		}

// 		if w.Body.String() != item.BodyResp {
// 			t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
// 				caseNum+1, w.Body.String(), item.BodyResp)
// 		}

// 		if item.testType == correctCase {
// 			if !testSessionDB.IsSessionByCookie(item.CookieReq.Value) {
// 				t.Errorf("TestCase [%d]:\nsession was not created", caseNum+1)
// 			}
// 			testSessionDB.DropCookies()
// 		}
// 	}
// }

// func TestSignup(t *testing.T) {
// 	t.Parallel()

// 	email := "testSignup1@mail.ru"
// 	password := "123456qQ"

// 	expectedUser := models.MakeUser(2, email, password)

// 	cases := []TestCase{
// 		{
// 			testType:   correctCase,
// 			BodyReq:    bytes.NewReader([]byte(`{"email":"` + email + `","password":"` + password + `"}`)),
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":200,"body":null}`,
// 		},
// 		{
// 			testType:   wrongCase,
// 			BodyReq:    bytes.NewReader([]byte(`wrong input data`)),
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":400,"body":null}`,
// 		},
// 		{
// 			testType:   wrongCase,
// 			BodyReq:    bytes.NewReader([]byte(`{"email":"firsUser@mail.ru","password":"EmailAlreadyExists"}`)),
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":1001,"body":null}`,
// 		},
// 	}

// 	testDB := _userRepo.NewMockDB()
// 	testSessionDB := session.NewSessionDB()

// 	timeoutContext := configs.Timeouts.ContextTimeout
// 	userHandler := &UserHandler{
// 		UserUCase: _userUCase.NewUserUsecase(testDB, testSessionDB, timeoutContext),
// 	}

// 	for caseNum, item := range cases {
// 		testDB.DropUsers(context.TODO())
// 		_, err := testDB.CreateUser(context.TODO(), models.LoginUser{
// 			Email:    "firsUser@mail.ru",
// 			Password: "123456qQ",
// 		})
// 		if err != nil {
// 			t.Errorf("Create user failed")
// 		}

// 		r := httptest.NewRequest("POST", "/api/v1/signup", item.BodyReq)
// 		r.AddCookie(&item.CookieReq)
// 		w := httptest.NewRecorder()

// 		userHandler.SignupHandler(w, r)

// 		if w.Code != item.StatusCode {
// 			t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
// 				caseNum+1, w.Code, item.StatusCode)
// 		}

// 		if w.Body.String() != item.BodyResp {
// 			t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
// 				caseNum+1, w.Body.String(), item.BodyResp)
// 		}

// 		if item.testType == correctCase {
// 			if !testSessionDB.IsSessionByCookie(item.CookieReq.Value) {
// 				t.Errorf("TestCase [%d]:\nsession was not created", caseNum+1)
// 			}
// 			testSessionDB.DropCookies()

// 			newUser, _ := testDB.GetUser(context.TODO(), email)
// 			if !reflect.DeepEqual(newUser, expectedUser) {
// 				t.Errorf("TestCase [%d]:\nuser was not created", caseNum+1)
// 			}
// 		}
// 	}
// }

// func TestLogout(t *testing.T) {
// 	t.Parallel()
// 	cases := []TestCase{
// 		{
// 			testType: correctCase,
// 			CookieReq: http.Cookie{
// 				Name:  "sessionId",
// 				Value: "123",
// 			},
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":200,"body":null}`,
// 		},
// 		{
// 			testType: wrongCase,
// 			CookieReq: http.Cookie{
// 				Name:  "sessionId",
// 				Value: "case wrong cookie",
// 			},
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":500,"body":null}`,
// 		},
// 		{
// 			testType:   wrongCase,
// 			CookieReq:  http.Cookie{},
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":404,"body":null}`,
// 		},
// 	}

// 	testDB := _userRepo.NewMockDB()
// 	testSessionDB := session.NewSessionDB()

// 	timeoutContext := configs.Timeouts.ContextTimeout
// 	sessionHandler := &_sessionDelivery.SessionHandler{
// 		UserUCase: _userUCase.NewUserUsecase(testDB, testSessionDB, timeoutContext),
// 	}

// 	user, err := testDB.CreateUser(context.TODO(), models.LoginUser{
// 		Email:    "testLogout1@mail.ru",
// 		Password: "123456qQ",
// 	})
// 	if err != nil {
// 		t.Errorf("Create user failed")
// 	}

// 	for caseNum, item := range cases {
// 		r := httptest.NewRequest("GET", "/api/v1/logout", item.BodyReq)
// 		r.AddCookie(&item.CookieReq)
// 		w := httptest.NewRecorder()

// 		err := testSessionDB.NewSessionCookie("123", user.ID)
// 		if err != nil {
// 			t.Errorf("New session Cookie error")
// 		}

// 		sessionHandler.LogoutHandler(w, r)

// 		if w.Code != item.StatusCode {
// 			t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
// 				caseNum+1, w.Code, item.StatusCode)
// 		}

// 		if w.Body.String() != item.BodyResp {
// 			t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
// 				caseNum+1, w.Body.String(), item.BodyResp)
// 		}

// 		if testSessionDB.IsSessionByCookie(item.CookieReq.Value) && item.testType == correctCase {
// 			t.Errorf("TestCase [%d]:\nuser session not ended", caseNum+1)
// 		}
// 	}
// }

// func TestNextUser(t *testing.T) {
// 	t.Parallel()
// 	cases := []TestCase{
// 		{
// 			testType: correctCase,
// 			BodyReq:  bytes.NewReader([]byte(`{"id":321}`)),
// 			CookieReq: http.Cookie{
// 				Name:  "sessionId",
// 				Value: "123",
// 			},
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":200,"body":{"id":1,"email":"testNextUser1@mail.ru"}}`,
// 		},
// 		{
// 			testType: wrongCase,
// 			BodyReq:  bytes.NewReader([]byte(`{"id":321}`)),
// 			CookieReq: http.Cookie{
// 				Name:  "sessionId",
// 				Value: "case wrong cookie",
// 			},
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":404,"body":null}`,
// 		},
// 		{
// 			testType:   wrongCase,
// 			BodyReq:    bytes.NewReader([]byte(`{"id":321}`)),
// 			CookieReq:  http.Cookie{},
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":404,"body":null}`,
// 		},
// 		{
// 			testType: wrongCase,
// 			BodyReq:  bytes.NewReader([]byte("wrong json")),
// 			CookieReq: http.Cookie{
// 				Name:  "sessionId",
// 				Value: "123",
// 			},
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":400,"body":null}`,
// 		},
// 		{
// 			testType: wrongCase,
// 			BodyReq:  bytes.NewReader([]byte(`{"id":1}`)),
// 			CookieReq: http.Cookie{
// 				Name:  "sessionId",
// 				Value: "123",
// 			},
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":404,"body":null}`,
// 		},
// 	}

// 	testDB := _userRepo.NewMockDB()
// 	testSessionDB := session.NewSessionDB()

// 	timeoutContext := configs.Timeouts.ContextTimeout
// 	userHandler := &UserHandler{
// 		UserUCase: _userUCase.NewUserUsecase(testDB, testSessionDB, timeoutContext),
// 	}

// 	_, err := testDB.CreateUser(context.TODO(), models.LoginUser{
// 		Email:    "testNextUser1@mail.ru",
// 		Password: "123456qQ\"",
// 	})
// 	if err != nil {
// 		t.Errorf("Create user failed")
// 	}

// 	currenUser, _ := testDB.CreateUser(context.TODO(), models.LoginUser{
// 		Email:    "testCurrUser1@mail.ru",
// 		Password: "123456qQ\"",
// 	})
// 	err = testSessionDB.NewSessionCookie("123", currenUser.ID)
// 	if err != nil {
// 		t.Error("New session error")
// 	}

// 	for caseNum, item := range cases {
// 		r := httptest.NewRequest("POST", "/api/v1/nextswipeuser", item.BodyReq)
// 		r.AddCookie(&item.CookieReq)
// 		w := httptest.NewRecorder()

// 		userHandler.NextUserHandler(w, r)

// 		if w.Code != item.StatusCode {
// 			t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
// 				caseNum+1, w.Code, item.StatusCode)
// 		}

// 		if w.Body.String() != item.BodyResp {
// 			t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
// 				caseNum+1, w.Body.String(), item.BodyResp)
// 		}

// 		resp, ok := testDB.IsSwiped(context.TODO(), currenUser.ID, 321)
// 		if ok != nil {
// 			t.Error("IsSwiped error")
// 		}
// 		if !resp && item.testType == correctCase {
// 			t.Errorf("TestCase [%d]:\nswipe not saved", caseNum+1)
// 		}
// 		testDB.DropSwipes(context.TODO())
// 	}
// }

// func TestEditProfile(t *testing.T) {
// 	t.Parallel()

// 	requestUser := models.User{
// 		Name:        "testEdit",
// 		Date:        "1999-10-25",
// 		Description: "Description Description Description Description",
// 		Imgs:        []string{"/img/testEdit/"},
// 		Tags:        []string{"Tags", "Tags", "Tags", "Tags", "Tags"},
// 	}
// 	bodyReq, err := json.Marshal(requestUser)
// 	if err != nil {
// 		t.Error("marshal error")
// 	}

// 	expectedUser := models.MakeUser(1, "testEdit@mail.ru", "123456qQ")
// 	err = expectedUser.FillProfile(&models.User{
// 		Name:        requestUser.Name,
// 		Date:        requestUser.Date,
// 		Description: requestUser.Description,
// 		Imgs:        requestUser.Imgs,
// 		Tags:        requestUser.Tags,
// 	})
// 	if err != nil {
// 		t.Error("fill profile error")
// 	}

// 	BodyRespByte, _ := json.Marshal(responses.JSON{
// 		Status: StatusOK,
// 		Body:   expectedUser,
// 	})

// 	cases := []TestCase{
// 		{
// 			testType: correctCase,
// 			BodyReq:  bytes.NewReader(bodyReq),
// 			CookieReq: http.Cookie{
// 				Name:  "sessionId",
// 				Value: "123",
// 			},
// 			StatusCode: http.StatusOK,
// 			BodyResp:   string(BodyRespByte),
// 		},
// 		{
// 			testType: wrongCase,
// 			BodyReq:  bytes.NewReader(bodyReq),
// 			CookieReq: http.Cookie{
// 				Name:  "sessionId",
// 				Value: "case wrong cookie",
// 			},
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":404,"body":null}`,
// 		},
// 		{
// 			testType:   wrongCase,
// 			BodyReq:    bytes.NewReader(bodyReq),
// 			CookieReq:  http.Cookie{},
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":404,"body":null}`,
// 		},
// 		{
// 			testType: wrongCase,
// 			BodyReq:  bytes.NewReader([]byte(`{"name":"testEdit","date":"wrong-format-data","description":"Description Description Description Description","imgSrc":"/img/testEdit/","tags":["Tags","Tags","Tags","Tags","Tags"]}`)),
// 			CookieReq: http.Cookie{
// 				Name:  "sessionId",
// 				Value: "123",
// 			},
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":400,"body":null}`,
// 		},
// 		{
// 			testType: wrongCase,
// 			BodyReq:  bytes.NewReader([]byte(`wrong data`)),
// 			CookieReq: http.Cookie{
// 				Name:  "sessionId",
// 				Value: "123",
// 			},
// 			StatusCode: http.StatusOK,
// 			BodyResp:   `{"status":400,"body":null}`,
// 		},
// 	}

// 	testDB := _userRepo.NewMockDB()
// 	testSessionDB := session.NewSessionDB()

// 	timeoutContext := configs.Timeouts.ContextTimeout
// 	userHandler := &UserHandler{
// 		UserUCase: _userUCase.NewUserUsecase(testDB, testSessionDB, timeoutContext),
// 	}

// 	for caseNum, item := range cases {
// 		testDB.DropUsers(context.TODO())
// 		testSessionDB.DropCookies()
// 		currenUser, err := testDB.CreateUser(context.TODO(), models.LoginUser{
// 			Email:    expectedUser.Email,
// 			Password: "123456qQ",
// 		})
// 		if err != nil {
// 			t.Error("create user error")
// 		}
// 		err = testSessionDB.NewSessionCookie("123", currenUser.ID)
// 		if err != nil {
// 			t.Error("new session error")
// 		}

// 		r := httptest.NewRequest("POST", "/api/v1/editprofile", item.BodyReq)
// 		r.AddCookie(&item.CookieReq)
// 		w := httptest.NewRecorder()

// 		userHandler.EditProfileHandler(w, r)

// 		if w.Code != item.StatusCode {
// 			t.Errorf("TestCase [%d]:\nwrongCase StatusCode: \ngot %d\nexpected %d",
// 				caseNum+1, w.Code, item.StatusCode)
// 		}

// 		if w.Body.String() != item.BodyResp {
// 			t.Errorf("TestCase [%d]:\nwrongCase Response: \ngot %s\nexpected %s",
// 				caseNum+1, w.Body.String(), item.BodyResp)
// 		}

// 		if item.testType == correctCase {
// 			updateUser, err := testDB.GetUser(context.TODO(), currenUser.Email)
// 			if err != nil {
// 				t.Errorf("TestCase [%d]:\nprofile was not created", caseNum+1)
// 			}
// 			if !reflect.DeepEqual(updateUser, expectedUser) {
// 				t.Errorf("TestCase [%d]:\nwrong profile: \ngot %v\nexpected %v",
// 					caseNum+1, updateUser, expectedUser)
// 			}
// 		}
// 	}
// }
