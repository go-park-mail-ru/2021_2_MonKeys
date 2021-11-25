package delivery

import (
	"dripapp/internal/microservices/chat/mocks"
	"dripapp/internal/microservices/chat/models"
	"dripapp/internal/pkg/logger"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
)

type TestCaseWS struct {
	MockWS []interface{}
	err    error
}

var (
	errInvalidSyntax = errors.New("invalid syntax")
)

func TestReadMessage(t *testing.T) {
	t.Parallel()

	mockConn := &mocks.WS{}
	call := mockConn.On("ReadJSON", mock.AnythingOfType("*models.Message"))

	msgWS := &MessagesWS{
		conn:   mockConn,
		logger: logger.DripLogger,
	}

	cases := []TestCaseWS{
		{
			MockWS: []interface{}{
				nil,
			},
			err: nil,
		},
		{
			MockWS: []interface{}{
				errInvalidSyntax,
			},
			err: errInvalidSyntax,
		},
	}

	for caseNum, item := range cases {
		call.Return(item.MockWS...)

		var message models.Message
		err := msgWS.ReadMessage(&message)

		if err != item.err {
			t.Errorf("TestCase [%d]:\nwrongCase: \ngot %d\nexpected %d",
				caseNum, err, item.err)
		}
	}
}

func TestWriteMessage(t *testing.T) {
	t.Parallel()

	mockConn := &mocks.WS{}
	call := mockConn.On("WriteJSON", mock.AnythingOfType("models.Message"))

	msgWS := &MessagesWS{
		conn:   mockConn,
		logger: logger.DripLogger,
	}

	cases := []TestCaseWS{
		{
			MockWS: []interface{}{
				nil,
			},
			err: nil,
		},
		{
			MockWS: []interface{}{
				errInvalidSyntax,
			},
			err: errInvalidSyntax,
		},
	}

	for caseNum, item := range cases {
		call.Return(item.MockWS...)

		var message models.Message
		err := msgWS.WriteMessage(message)

		if err != item.err {
			t.Errorf("TestCase [%d]:\nwrongCase: \ngot %d\nexpected %d",
				caseNum, err, item.err)
		}
	}
}
