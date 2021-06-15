package sharesies_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"testing"

	"github.com/deividfortuna/sharesies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDoType func(req *http.Request) (*http.Response, error)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()

	args := m.Called(req.Method, req.URL, string(body))

	return args.Get(0).(*http.Response), nil
}

func Test_New(t *testing.T) {
	s, err := sharesies.New(nil)

	assert.Nil(t, err)
	assert.NotNil(t, s)
}

func Test_New_HttpClient(t *testing.T) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	s, err := sharesies.New(client)

	assert.Nil(t, err)
	assert.NotNil(t, s)
}

func Test_New_HttpClient_MissConfigurated(t *testing.T) {
	client := http.DefaultClient

	s, err := sharesies.New(client)

	assert.Nil(t, s)
	assert.Equal(t, sharesies.ErrNoJarDefine, err)
}

func Test_Authenticate(t *testing.T) {
	mockClient := &MockClient{}
	authSuccess(mockClient)

	s := sharesies.Sharesies{
		HttpClient: mockClient,
	}

	ctx := context.Background()
	r, err := s.Authenticate(ctx, &sharesies.Credentials{Username: "username", Password: "password"})

	mockClient.AssertExpectations(t)

	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.True(t, r.Authenticated)
}

func Test_Authenticate_Fail(t *testing.T) {
	mockClient := &MockClient{}
	url, _ := url.Parse("https://app.sharesies.nz/api/identity/login")
	body := marshal(&sharesies.Map{"email": "username", "password": "password", "remember": true})

	mockClient.On("Do", http.MethodPost, url, body).Return(&http.Response{StatusCode: http.StatusUnauthorized}, nil)

	s := sharesies.Sharesies{
		HttpClient: mockClient,
	}

	ctx := context.Background()
	r, err := s.Authenticate(ctx, &sharesies.Credentials{Username: "username", Password: "password"})

	mockClient.AssertExpectations(t)

	assert.Equal(t, sharesies.ErrHttpRequest, err)
	assert.Nil(t, r)
}

func Test_Profile(t *testing.T) {
	mockClient := &MockClient{}
	authSuccess(mockClient)

	profileUrl, _ := url.Parse("https://app.sharesies.nz/api/identity/check")
	profileBody, _ := os.Open("testdata/profile.json")

	mockClient.On("Do", http.MethodGet, profileUrl, mock.Anything).Return(&http.Response{StatusCode: http.StatusOK, Body: profileBody}, nil)

	s := sharesies.Sharesies{
		HttpClient: mockClient,
	}

	ctx := context.Background()
	s.Authenticate(ctx, &sharesies.Credentials{Username: "username", Password: "password"})

	p, err := s.Profile(ctx)
	mockClient.AssertExpectations(t)

	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.True(t, p.Authenticated)
}

func Test_Instruments(t *testing.T) {
	mockClient := &MockClient{}
	authSuccess(mockClient)

	instrumentRequest := &sharesies.InstrumentsRequest{
		Page:            1,
		Perpage:         60,
		Sort:            "relevance",
		Pricechangetime: "1y",
		Query:           "apple",
	}

	instrumentsUrl, _ := url.Parse("https://data.sharesies.nz/api/v1/instruments")
	instrumentsBody, _ := os.Open("testdata/instruments.json")

	mockClient.On("Do", http.MethodPost, instrumentsUrl, marshal(instrumentRequest)).Return(&http.Response{StatusCode: http.StatusOK, Body: instrumentsBody}, nil)

	s := sharesies.Sharesies{
		HttpClient: mockClient,
	}

	ctx := context.Background()
	s.Authenticate(ctx, &sharesies.Credentials{Username: "username", Password: "password"})

	i, err := s.Instruments(ctx, instrumentRequest)
	mockClient.AssertExpectations(t)

	assert.Nil(t, err)
	assert.NotNil(t, i)
}

func Test_CostBuy(t *testing.T) {
	mockClient := &MockClient{}
	authSuccess(mockClient)
	reAuthSuccess(mockClient)

	costBuyUrl, _ := url.Parse("https://app.sharesies.nz/api/order/cost-buy")
	costBuyBody, _ := os.Open("testdata/costbuy.json")
	body := marshal(&sharesies.CostBuyRequest{
		FundID:     "b8b7ef58-b270-4762-a256-9d68aebc3e23",
		ActingAsID: "USER_ID",
		Order: &sharesies.Order{
			Type:           sharesies.OrderTypeDollarMarket,
			CurrencyAmount: "10.00",
		},
	})

	mockClient.On("Do", http.MethodPost, costBuyUrl, body).Return(&http.Response{StatusCode: http.StatusOK, Body: costBuyBody}, nil)

	s := sharesies.Sharesies{
		HttpClient: mockClient,
	}

	ctx := context.Background()
	s.Authenticate(ctx, &sharesies.Credentials{Username: "username", Password: "password"})

	i, err := s.CostBuy(ctx, "b8b7ef58-b270-4762-a256-9d68aebc3e23", 10.00)
	mockClient.AssertExpectations(t)

	assert.Nil(t, err)
	assert.NotNil(t, i)
}

func reAuthSuccess(mockClient *MockClient) {
	reAuthUrl, _ := url.Parse("https://app.sharesies.nz/api/identity/reauthenticate")
	authBody, _ := os.Open("testdata/authenticated.json")
	body := marshal(&sharesies.Map{"password": "password", "acting_as_id": "USER_ID"})

	mockClient.On("Do", http.MethodPost, reAuthUrl, body).Return(&http.Response{StatusCode: http.StatusOK, Body: authBody}, nil)
}

func authSuccess(mockClient *MockClient) {
	authUrl, _ := url.Parse("https://app.sharesies.nz/api/identity/login")
	authBody, _ := os.Open("testdata/authenticated.json")
	body := marshal(&sharesies.Map{"email": "username", "password": "password", "remember": true})

	mockClient.On("Do", http.MethodPost, authUrl, body).Return(&http.Response{StatusCode: http.StatusOK, Body: authBody}, nil)
}

func marshal(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
