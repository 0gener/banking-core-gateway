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

type createAccountTestCase struct {
	name                   string
	request                *createAccountRequest
	expectedStatusCode     int
	expectedResponse       *createAccountResponse
	expectedError          error
	accountsClientResponse *proto.CreateAccountResponse
	accountsClientError    error
}

type getAccountTestCase struct {
	name                   string
	expectedStatusCode     int
	expectedResponse       *getAccountResponse
	expectedError          error
	accountsClientResponse *proto.GetAccountResponse
	accountsClientError    error
}

type testAccountsClient struct {
	t                     *testing.T
	createAccountTestCase createAccountTestCase
	getAccountTestCase    getAccountTestCase
}

func (c *testAccountsClient) CreateAccount(ctx context.Context, in *proto.CreateAccountRequest, opts ...grpc.CallOption) (*proto.CreateAccountResponse, error) {
	if *c.createAccountTestCase.request.Currency != in.Currency {
		c.t.Errorf("expected currency %s, got %s", *c.createAccountTestCase.request.Currency, in.Currency)
	}

	return c.createAccountTestCase.accountsClientResponse, c.createAccountTestCase.accountsClientError
}

func (c *testAccountsClient) GetAccount(ctx context.Context, in *proto.GetAccountRequest, opts ...grpc.CallOption) (*proto.GetAccountResponse, error) {
	return c.getAccountTestCase.accountsClientResponse, c.getAccountTestCase.accountsClientError
}

func TestCreateAccountHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	currency := "1234"
	testCases := []createAccountTestCase{
		{
			name: "valid_request",
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
			name:                   "no_request",
			request:                nil,
			expectedStatusCode:     http.StatusBadRequest,
			expectedResponse:       nil,
			expectedError:          nil,
			accountsClientResponse: nil,
			accountsClientError:    nil,
		},
		{
			name:                   "request_without_currency",
			request:                &createAccountRequest{},
			expectedStatusCode:     http.StatusBadRequest,
			expectedResponse:       nil,
			expectedError:          nil,
			accountsClientResponse: nil,
			accountsClientError:    nil,
		},
		{
			name: "accounts_service_error",
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

			accountsClient := testAccountsClient{
				t:                     t,
				createAccountTestCase: tc,
			}

			rec := httptest.NewRecorder()
			r := New(&testJwtMiddleware{}, &accountsClient)

			var reader io.Reader
			if tc.request != nil {
				body, _ := json.Marshal(tc.request)
				reader = bytes.NewBuffer(body)
			}
			request := httptest.NewRequest("POST", "/accounts", reader)

			r.ServeHTTP(rec, request)

			if rec.Code != tc.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tc.expectedStatusCode, rec.Code)
			}

			bytes, _ := ioutil.ReadAll(rec.Body)
			var response *createAccountResponse

			if len(bytes) > 0 {
				response = &createAccountResponse{}
				json.Unmarshal(bytes, response)
			}

			if response == nil || tc.expectedResponse == nil {
				if response != tc.expectedResponse {
					t.Errorf("expected response %s, got %s", tc.expectedResponse, response)
				}
			} else if *response != *tc.expectedResponse {
				t.Errorf("expected response %s, got %s", *tc.expectedResponse, *response)
			}
		})
	}
}

func TestGetAccountHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testCases := []getAccountTestCase{
		{
			name:               "valid_request",
			expectedStatusCode: http.StatusOK,
			expectedResponse: &getAccountResponse{
				AccountNumber: "1234",
				Currency:      "EUR",
			},
			expectedError: nil,
			accountsClientResponse: &proto.GetAccountResponse{
				Account: &proto.Account{
					AccountNumber: "1234",
					Currency:      "EUR",
				},
			},
			accountsClientError: nil,
		},
		{
			name:                   "accounts_service_error",
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

			accountsClient := testAccountsClient{
				t:                  t,
				getAccountTestCase: tc,
			}

			rec := httptest.NewRecorder()
			r := New(&testJwtMiddleware{}, &accountsClient)

			request := httptest.NewRequest("GET", "/accounts", nil)

			r.ServeHTTP(rec, request)

			if rec.Code != tc.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tc.expectedStatusCode, rec.Code)
			}

			bytes, _ := ioutil.ReadAll(rec.Body)
			var response *getAccountResponse

			if len(bytes) > 0 {
				response = &getAccountResponse{}
				json.Unmarshal(bytes, response)
			}

			if response == nil || tc.expectedResponse == nil {
				if response != tc.expectedResponse {
					t.Errorf("expected response %v, got %v", tc.expectedResponse, response)
				}
			} else if *response != *tc.expectedResponse {
				t.Errorf("expected response %v, got %v", *tc.expectedResponse, *response)
			}
		})
	}
}
