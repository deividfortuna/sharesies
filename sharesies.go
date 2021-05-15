package sharesies

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const (
	endpointIdentityLogin  = "https://app.sharesies.nz/api/identity/login"
	endpointIdentityCheck  = "https://app.sharesies.nz/api/identity/check"
	endpointIdentityReAuth = "https://app.sharesies.nz/api/identity/reauthenticate"
	endpointInstruments    = "https://data.sharesies.nz/api/v1/instruments"
	endpointCostBuy        = "https://app.sharesies.nz/api/order/cost-buy"
	endpointCreateBuy      = "https://app.sharesies.nz/api/order/create-buy"
)

type Map map[string]interface{}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Sharesies struct {
	HttpClient HTTPClient
	creds      *SharesiesCredentials
	session    *sharesiesSession
}

// Sharesies Cedentials
type SharesiesCredentials struct {
	Username string
	Password string
}

type sharesiesRequest struct {
	url      string
	method   string
	body     interface{}
	response interface{}
	headers  map[string]string
}

type sharesiesSession struct {
	token   *jwt.Token
	profile *ProfileResponse
}

// New credentials struct
func NewCredentials(username string, passoword string) *SharesiesCredentials {
	return &SharesiesCredentials{Username: username, Password: passoword}
}

// NewClient returns a new Sharesies Client instance
func NewClient() (*Sharesies, error) {
	j, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Jar: j,
	}

	return &Sharesies{
		httpClient,
		nil,
		nil,
	}, nil
}

func (s *Sharesies) Authenticate(ctx context.Context, creds *SharesiesCredentials) (*ProfileResponse, error) {
	p := &ProfileResponse{}
	req := &sharesiesRequest{
		method:   http.MethodPost,
		url:      endpointIdentityLogin,
		body:     &Map{"email": creds.Username, "password": creds.Password, "remember": true},
		response: p,
	}
	err := s.request(ctx, req)
	if !p.Authenticated {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	errCtx := s.authenticated(p)
	if errCtx != nil {
		return nil, errCtx
	}

	s.creds = creds

	return p, err
}

// Return Sharesies Profile
func (s *Sharesies) Profile(ctx context.Context) (*ProfileResponse, error) {
	p := &ProfileResponse{}
	req := &sharesiesRequest{
		method:   http.MethodGet,
		url:      endpointIdentityCheck,
		response: p,
	}
	err := s.request(ctx, req)
	if err != nil {
		return nil, err
	}

	er := s.authenticated(p)
	if er != nil {
		return nil, er
	}

	return p, nil
}

// Return Companies/Funds listed on Sharesies
func (s *Sharesies) Instruments(ctx context.Context, request *InstrumentsRequest) (*InstrumentResponse, error) {
	r := &InstrumentResponse{}
	h, errH := s.headers(ctx)
	if errH != nil {
		return nil, errH
	}

	req := &sharesiesRequest{
		method:   http.MethodPost,
		url:      endpointInstruments,
		body:     request,
		response: r, headers: h,
	}
	err := s.request(ctx, req)
	return r, err
}

// Cost to buy stocks from the NZX Market
func (s *Sharesies) CostBuy(ctx context.Context, fundId string, amount float64) (*CostBuyResponse, error) {
	r := &CostBuyResponse{}
	o := &Order{Type: OrderTypeDollarMarket, CurrencyAmount: fmt.Sprintf("%.2f", amount)}
	cr := &CostBuyRequest{
		FundID:     fundId,
		ActingAsID: s.session.profile.UserList[0].ID,
		Order:      o,
	}

	s.reAuthenticate(ctx)

	req := &sharesiesRequest{
		method:   http.MethodPost,
		url:      endpointCostBuy,
		body:     cr,
		response: r,
	}

	err := s.request(ctx, req)
	return r, err
}

// Purchase stocks from the NZX Market
func (s *Sharesies) Buy(ctx context.Context, costBuy *CostBuyResponse) (*ProfileResponse, error) {
	r := &ProfileResponse{}

	br := &CreateBuyRequest{
		FundID:           costBuy.FundID,
		ActingAsID:       s.session.profile.UserList[0].ID,
		Order:            costBuy.Request,
		PaymentBreakdown: &costBuy.PaymentBreakdown,
		IdempotencyKey:   uuid.NewString(),
		ExpectedFee:      costBuy.ExpectedFee,
	}

	s.reAuthenticate(ctx)

	req := &sharesiesRequest{
		method:   http.MethodPost,
		url:      endpointCreateBuy,
		body:     br,
		response: r,
	}

	err := s.request(ctx, req)
	return r, err
}

func (s *Sharesies) authenticated(p *ProfileResponse) error {
	claims := &jwt.StandardClaims{}
	token, _, err := new(jwt.Parser).ParseUnverified(p.DistillToken, claims)
	if err != nil {
		return nil
	}

	s.session = &sharesiesSession{
		token:   token,
		profile: p,
	}

	return nil
}

func (s *Sharesies) headers(ctx context.Context) (map[string]string, error) {
	err := s.session.token.Claims.Valid()
	if err != nil {
		_, err := s.reAuthenticate(ctx)
		if err != nil {
			return nil, err
		}
	}

	return map[string]string{
		"Authorization": "Bearer " + s.session.token.Raw,
	}, nil
}

func (s *Sharesies) reAuthenticate(ctx context.Context) (*ProfileResponse, error) {
	p := &ProfileResponse{}
	req := &sharesiesRequest{
		method:   http.MethodPost,
		url:      endpointIdentityReAuth,
		body:     &Map{"password": s.creds.Password, "acting_as_id": s.session.profile.UserList[0].ID},
		response: p,
	}
	err := s.request(ctx, req)
	if err != nil {
		return nil, err
	}

	if !p.Authenticated {
		return nil, fmt.Errorf("failed to re-authenticate: %w", err)
	}

	e := s.authenticated(p)
	if e != nil {
		return nil, e
	}

	return p, nil
}

func (s *Sharesies) request(ctx context.Context, request *sharesiesRequest) error {
	b := &bytes.Buffer{}
	e := json.NewEncoder(b)
	e.Encode(request.body)

	req, err := http.NewRequestWithContext(ctx, request.method, request.url, b)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 Firefox/71.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/json")

	if request.headers != nil {
		for key, value := range request.headers {
			req.Header.Add(key, value)
		}
	}

	res, err := s.HttpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("expected 200 code, received: %v", res.StatusCode)
	}

	bd, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	json.Unmarshal(bd, &request.response)
	return nil
}
