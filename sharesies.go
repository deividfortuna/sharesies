package sharesies

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	endpointCostSell       = "https://app.sharesies.nz/api/order/cost-sell"
	endpointCreateSell     = "https://app.sharesies.nz/api/order/create-sell"
)

type Map map[string]interface{}

var ErrNoJarDefine = errors.New("HttpClient must have a cookie jar defined")
var ErrAuthentication = errors.New("authentication failed")
var ErrHttpRequest = errors.New("request to sharesies failed")

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Sharesies struct {
	HttpClient HTTPClient
	creds      *Credentials
	session    *tokenSession
}

// Credentials Sharesies
type Credentials struct {
	Username string
	Password string
}

type tokenSession struct {
	token   *jwt.Token
	profile *ProfileResponse
}

// New returns a new Sharesies Client instance
func New(client *http.Client) (*Sharesies, error) {
	if client == nil {
		j, err := cookiejar.New(nil)
		if err != nil {
			return nil, err
		}

		client = &http.Client{
			Jar: j,
		}
	}

	if client.Jar == nil {
		return nil, ErrNoJarDefine
	}

	return &Sharesies{
		client,
		nil,
		nil,
	}, nil
}

func (s *Sharesies) Authenticate(ctx context.Context, creds *Credentials) (*ProfileResponse, error) {
	p := &ProfileResponse{}
	body := &Map{"email": creds.Username, "password": creds.Password, "remember": true}

	err := s.request(ctx, http.MethodPost, nil, endpointIdentityLogin, body, p)
	if err != nil {
		return nil, err
	}

	if !p.Authenticated {
		return nil, ErrAuthentication
	}

	err = s.authenticated(p)
	if err != nil {
		return nil, err
	}

	s.creds = creds

	return p, err
}

// Profile returns Sharesies Profile
func (s *Sharesies) Profile(ctx context.Context) (*ProfileResponse, error) {
	p := &ProfileResponse{}

	err := s.request(ctx, http.MethodGet, nil, endpointIdentityCheck, nil, p)
	if err != nil {
		return nil, err
	}

	err = s.authenticated(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Instruments returns Companies/Funds listed on Sharesies
func (s *Sharesies) Instruments(ctx context.Context, request *InstrumentsRequest) (*InstrumentResponse, error) {
	r := &InstrumentResponse{}
	h, err := s.headers(ctx)
	if err != nil {
		return nil, err
	}

	err = s.request(ctx, http.MethodPost, h, endpointInstruments, request, r)
	return r, err
}

// CostBuy return Cost to buy stocks from the NZX Market
func (s *Sharesies) CostBuy(ctx context.Context, fundId string, amount float64) (*CostBuyResponse, error) {
	r := &CostBuyResponse{}
	o := &OrderBuy{Type: OrderTypeDollarMarket, CurrencyAmount: fmt.Sprintf("%.2f", amount)}
	cr := &CostBuyRequest{
		FundID:     fundId,
		ActingAsID: s.session.profile.UserList[0].ID,
		Order:      o,
	}

	s.reAuthenticate(ctx)

	err := s.request(ctx, http.MethodPost, nil, endpointCostBuy, cr, r)
	return r, err
}

// Buy purchase stocks from the NZX Market
func (s *Sharesies) Buy(ctx context.Context, costBuy *CostBuyResponse) (*ProfileResponse, error) {
	r := &ProfileResponse{}

	br := &CreateBuyRequest{
		FundID:           costBuy.FundID,
		ActingAsID:       s.session.profile.UserList[0].ID,
		Order:            costBuy.Request,
		PaymentBreakdown: costBuy.PaymentBreakdown,
		IdempotencyKey:   uuid.NewString(),
		ExpectedFee:      costBuy.ExpectedFee,
	}

	s.reAuthenticate(ctx)

	err := s.request(ctx, http.MethodPost, nil, endpointCreateBuy, br, r)
	return r, err
}

func (s *Sharesies) CostSell(ctx context.Context, foundId string, shareAmount float64) (*CostSellResponse, error) {
	r := &CostSellResponse{}
	o := &OrderSell{Type: OrderTypeShareMarket, ShareAmount: fmt.Sprintf("%.6f", shareAmount)}
	sr := &CostSellRequest{FundID: foundId, ActingAsID: s.session.profile.UserList[0].ID, Order: o}

	_, err := s.reAuthenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = s.request(ctx, http.MethodPost, nil, endpointCostSell, sr, r)
	return r, err
}

func (s *Sharesies) Sell(ctx context.Context, sellBuy *CostSellResponse) (*ProfileResponse, error) {
	r := &ProfileResponse{}

	sr := CreateSellRequest{
		FundID: sellBuy.FundID,
		ActingAsID: s.session.profile.UserList[0].ID,
		Order: sellBuy.Request,
	}

	_, err := s.reAuthenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = s.request(ctx, http.MethodPost, nil, endpointCreateSell, sr, r)
	return r, err
}

func (s *Sharesies) authenticated(p *ProfileResponse) error {
	claims := &jwt.StandardClaims{}
	token, _, err := new(jwt.Parser).ParseUnverified(p.DistillToken, claims)
	if err != nil {
		return nil
	}

	s.session = &tokenSession{
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
	body := &Map{"password": s.creds.Password, "acting_as_id": s.session.profile.UserList[0].ID}

	err := s.request(ctx, http.MethodPost, nil, endpointIdentityReAuth, body, p)
	if err != nil {
		return nil, err
	}

	if !p.Authenticated {
		return nil, ErrAuthentication
	}

	e := s.authenticated(p)
	if e != nil {
		return nil, e
	}

	return p, nil
}

func (s *Sharesies) request(ctx context.Context, method string, headers map[string]string, url string, body interface{}, response interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 Firefox/71.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	res, err := s.HttpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return ErrHttpRequest
	}

	bd, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	return json.Unmarshal(bd, &response)
}
