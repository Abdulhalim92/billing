package handler

import (
	"billing/internal/model"
	"billing/internal/service"
	mock_service "billing/internal/service/mocks"
	"billing/logging"
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_CreateAccount(t *testing.T) {
	type mockBehavior func(r *mock_service.MockAccount, account *model.Account)

	tests := []struct {
		name                 string
		inputBody            string
		inputAccount         *model.Account
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"number": "123456", "balance": 100.00, "status": "active"}`,
			inputAccount: &model.Account{
				Number:  "123456",
				Balance: 100.00,
				Status:  "active",
			},
			mockBehavior: func(r *mock_service.MockAccount, account *model.Account) {
				r.EXPECT().CreateAccount(account).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"ID":1,"message":"Successfully create account"}`,
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"number": "123456", "balance": 100.00}`,
			inputAccount:         &model.Account{},
			mockBehavior:         func(r *mock_service.MockAccount, account *model.Account) {},
			expectedStatusCode:   400,
			expectedResponseBody: `"please send valid data"`,
		},
		{
			name:      "Service Error",
			inputBody: `{"number": "123456", "balance": 100.00, "status": "active"}`,
			inputAccount: &model.Account{
				Number:  "123456",
				Balance: 100.00,
				Status:  "active",
			},
			mockBehavior: func(r *mock_service.MockAccount, account *model.Account) {
				r.EXPECT().CreateAccount(account).Return(0, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `"something went wrong"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockAccount(c)
			test.mockBehavior(repo, test.inputAccount)

			r := gin.Default()
			l := logging.GetLogger()

			services := &service.Service{Account: repo}
			//services := &service.Service{Authorization:service.NewAuthService()}

			handler := Handler{r, services, l}
			handler.Engine.POST("/v1/create", handler.CreateAccount)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/v1/create", bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)

		})
	}
}
