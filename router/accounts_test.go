package router

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0gener/banking-core-accounts/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type testCase struct {
	name                   string
	request                *createAccountRequest
	expectedStatusCode     int
	expectedResponse       *createAccountResponse
	expectedError          error
	accountsClientResponse interface{}
	accountsClientError    error
}

type testAccountsClient struct {
	t  *testing.T
	tc testCase
}

func (c *testAccountsClient) CreateAccount(ctx context.Context, in *proto.CreateAccountRequest, opts ...grpc.CallOption) (*proto.CreateAccountResponse, error) {
	if *c.tc.request.Currency != in.Currency {
		c.t.Errorf("expected currency %s, got %s", *c.tc.request.Currency, in.Currency)
	}

	var response *proto.CreateAccountResponse
	if c.tc.accountsClientResponse != nil {
		response = c.tc.accountsClientResponse.(*proto.CreateAccountResponse)
	}

	return response, c.tc.accountsClientError
}

func (c *testAccountsClient) GetAccount(ctx context.Context, in *proto.GetAccountRequest, opts ...grpc.CallOption) (*proto.GetAccountResponse, error) {
	return c.tc.accountsClientResponse.(*proto.GetAccountResponse), c.tc.accountsClientError
}

func TestCreateAccountHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	currency := "1234"
	testCases := []testCase{
		{
			name: "valid request",
			request: &createAccountRequest{
				Currency: &currency,
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: &createAccountResponse{
				AccountNumber: "1234",
				Currency:      currency,
			},
			expectedError: nil,
			accountsClientResponse: &proto.CreateAccountResponse{
				Account: &proto.Account{
					AccountNumber: "1234",
					Currency:      currency,
				},
			},
			accountsClientError: nil,
		},
		{
			name:                   "no request",
			request:                nil,
			expectedStatusCode:     http.StatusBadRequest,
			expectedResponse:       nil,
			expectedError:          nil,
			accountsClientResponse: nil,
			accountsClientError:    nil,
		},
		{
			name:                   "request without currency",
			request:                &createAccountRequest{},
			expectedStatusCode:     http.StatusBadRequest,
			expectedResponse:       nil,
			expectedError:          nil,
			accountsClientResponse: nil,
			accountsClientError:    nil,
		},
		{
			name: "accounts service error",
			request: &createAccountRequest{
				Currency: &currency,
			},
			expectedStatusCode:     http.StatusInternalServerError,
			expectedResponse:       nil,
			expectedError:          nil,
			accountsClientResponse: nil,
			accountsClientError:    errors.New("just a random error"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			accountsClient := testAccountsClient{t, tc}
			accountsController := newAccountsController(&accountsClient)

			rec := httptest.NewRecorder()
			_, r := gin.CreateTestContext(rec)

			r.POST("/accounts", accountsController.createAccountHandler)

			var reader io.Reader
			if tc.request != nil {
				body, _ := json.Marshal(tc.request)
				reader = bytes.NewBuffer(body)
			}
			request := httptest.NewRequest("POST", "/accounts", reader)

			r.ServeHTTP(rec, request)

			if rec.Code != tc.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", rec.Code, tc.expectedStatusCode)
			}

			bytes, _ := ioutil.ReadAll(rec.Body)
			var response *createAccountResponse

			if len(bytes) > 0 {
				response = &createAccountResponse{}
				json.Unmarshal(bytes, response)
			}

			if response == nil || tc.expectedResponse == nil {
				if response != tc.expectedResponse {
					t.Errorf("expected response %s, got %s", response, tc.expectedResponse)
				}
			} else if *response != *tc.expectedResponse {
				t.Errorf("expected response %s, got %s", *response, *tc.expectedResponse)
			}
		})
	}
}
