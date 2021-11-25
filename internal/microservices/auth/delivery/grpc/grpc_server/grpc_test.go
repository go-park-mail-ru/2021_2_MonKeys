package grpcserver

// import (
// 	"context"
// 	_userModels "dripapp/internal/dripapp/models"
// 	_userMocks "dripapp/internal/dripapp/user/mocks"
// 	_authClient "dripapp/internal/microservices/auth/delivery/grpc/client"
// 	_authMocks "dripapp/internal/microservices/auth/mocks"
// 	_authModels "dripapp/internal/microservices/auth/models"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"google.golang.org/grpc"
// )

// func TestGrpcServer(t *testing.T) {
// 	type CheckStructInput struct {
// 		Session _authModels.Session
// 	}

// 	type CheckStructOutput struct {
// 		User _userModels.User
// 		Err  error
// 	}
// 	type AdditionalInfo struct {
// 		Anything error
// 	}

// 	type testCaseStruct struct {
// 		InputData      CheckStructInput
// 		OutputData     CheckStructOutput
// 		AdditionalInfo AdditionalInfo
// 	}

// 	testUser := _userModels.User{
// 		ID:          1,
// 		Email:       "valid@valid.ru",
// 		Password:    "!Nagdimaev2001",
// 		Name:        "Ilyagu",
// 		Gender:      "male",
// 		FromAge:     18,
// 		ToAge:       60,
// 		Date:        "2001-06-29",
// 		Age:         20,
// 		Description: "всем привет",
// 		Imgs:        []string{"img1", "img2"},
// 		Tags:        []string{"tag1", "tag2"},
// 	}

// 	testCases := []testCaseStruct{
// 		//test Ok
// 		{
// 			InputData: CheckStructInput{
// 				Session: _authModels.Session{
// 					Cookie: "lol",
// 					UserID: 1,
// 				},
// 			},
// 			OutputData: CheckStructOutput{
// 				User: testUser,
// 				Err:  nil,
// 			},
// 		},
// 	}

// 	authMockRepo := new(_authMocks.SessionRepository)
// 	userMockRepo := new(_userMocks.UserRepository)
// 	urlForTests := "127.0.0.1:8003"
// 	go StartAuthGrpcServer(authMockRepo, userMockRepo, urlForTests)
// 	grpcConn, err := grpc.Dial(urlForTests, grpc.WithInsecure())
// 	assert.Nil(t, err, "no error when start grpc conn")
// 	custGrpcClient := _authClient.NewAuthClient(grpcConn)

// 	for _, testCase := range testCases {
// 		userMockRepo.On("GetUserByID", context.Background(), testCase.InputData.Session.UserID).Return(testCase.OutputData.User, testCase.OutputData.Err)
// 		res, err := custGrpcClient.GetById(context.Background(), testCase.InputData.Session)
// 		assert.Equal(t, testCase.OutputData.Err, err)
// 		assert.Equal(t, testCase.OutputData.User, res)
// 	}

// 	type CheckStructInput2 struct {
// 		Cookie string
// 	}

// 	type CheckStructOutput2 struct {
// 		Session _authModels.Session
// 		Err     error
// 	}
// 	type AdditionalInfo2 struct {
// 		Anything error
// 	}

// 	type testCaseStruct2 struct {
// 		InputData      CheckStructInput2
// 		OutputData     CheckStructOutput2
// 		AdditionalInfo AdditionalInfo2
// 	}

// 	testCases2 := []testCaseStruct2{
// 		//test Ok
// 		{
// 			InputData: CheckStructInput2{
// 				Cookie: "lol",
// 			},
// 			OutputData: CheckStructOutput2{
// 				Session: _authModels.Session{
// 					Cookie: "lol",
// 					UserID: 1,
// 				},
// 				Err: nil,
// 			},
// 		},
// 	}

// 	for _, testCase2 := range testCases2 {
// 		authMockRepo.On("GetSessionByCookie", testCase2.InputData.Cookie).Return(testCase2.OutputData.Session, testCase2.OutputData.Err)
// 		res, err := custGrpcClient.GetFromSession(context.Background(), testCase2.InputData.Cookie)
// 		assert.Equal(t, testCase2.OutputData.Err, err)
// 		assert.Equal(t, testCase2.OutputData.Session, res)
// 	}
// }
