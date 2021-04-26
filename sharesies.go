package sharesies

import (
	"bytes"
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

type Sharesies struct {
	httpClient *http.Client
	creds      *SharesiesCredentials
	ctx        *sharesiesCtx
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

type sharesiesCtx struct {
	token   *jwt.Token
	profile *ProfileResponse
}

// New returns a new Sharesies Client instance
func New(creds *SharesiesCredentials) (*Sharesies, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Jar: jar,
	}

	i := &Sharesies{
		httpClient,
		creds,
		&sharesiesCtx{},
	}

	errCtx := i.sharesiesCtx(i.ctx)
	if errCtx != nil {
		return nil, errCtx
	}

	return i, nil
}

// Return Sharesies Profile
func (s *Sharesies) Profile() (*ProfileResponse, error) {
	p := &ProfileResponse{}
	req := &sharesiesRequest{
		method:   http.MethodGet,
		url:      endpointIdentityCheck,
		response: p,
	}
	err := s.request(req)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Return Companies/Funds listed on Sharesies
func (s *Sharesies) Instruments(request *InstrumentsRequest) (*InstrumentResponse, error) {
	r := &InstrumentResponse{}
	h, errH := s.getAuthHeaders()
	if errH != nil {
		return nil, errH
	}

	req := &sharesiesRequest{
		method:   http.MethodPost,
		url:      endpointInstruments,
		body:     request,
		response: r, headers: h,
	}
	err := s.request(req)
	return r, err
}

// Cost to buy stocks from the NZX Market
func (s *Sharesies) CostBuy(fundId string, amount float64) (*CostBuyResponse, error) {
	r := &CostBuyResponse{}
	o := &Order{Type: OrderTypeDollarMarket, CurrencyAmount: fmt.Sprintf("%.2f", amount)}
	cr := &CostBuyRequest{
		FundID:     fundId,
		ActingAsID: s.ctx.profile.UserList[0].ID,
		Order:      o,
	}

	s.reAuthenticate()

	req := &sharesiesRequest{
		method:   http.MethodPost,
		url:      endpointCostBuy,
		body:     cr,
		response: r,
	}

	err := s.request(req)
	return r, err
}

// Purchase stocks from the NZX Market
func (s *Sharesies) Buy(costBuy *CostBuyResponse) (*ProfileResponse, error) {
	r := &ProfileResponse{}

	br := &CreateBuyRequest{
		FundID:           costBuy.FundID,
		ActingAsID:       s.ctx.profile.UserList[0].ID,
		Order:            costBuy.Request,
		PaymentBreakdown: &costBuy.PaymentBreakdown,
		IdempotencyKey:   uuid.NewString(),
		ExpectedFee:      costBuy.ExpectedFee,
	}

	s.reAuthenticate()

	req := &sharesiesRequest{
		method:   http.MethodPost,
		url:      endpointCreateBuy,
		body:     br,
		response: r,
	}

	err := s.request(req)
	return r, err
}

func (s *Sharesies) sharesiesCtx(ctx *sharesiesCtx) error {
	p, err := s.authenticate(s.creds)
	if err != nil {
		return err
	}

	claims := &jwt.StandardClaims{}
	token, _, err := new(jwt.Parser).ParseUnverified(p.DistillToken, claims)
	if err != nil {
		return nil
	}

	ctx.token = token
	ctx.profile = p

	return nil
}

func (s *Sharesies) getAuthHeaders() (map[string]string, error) {
	err := s.ctx.token.Claims.Valid()
	if err != nil {
		_, err := s.reAuthenticate()
		if err != nil {
			return nil, err
		}
	}

	return map[string]string{
		"Authorization": "Bearer" + s.ctx.token.Raw,
	}, nil
}

func (s *Sharesies) authenticate(creds *SharesiesCredentials) (*ProfileResponse, error) {
	p := &ProfileResponse{}
	req := &sharesiesRequest{
		method:   http.MethodPost,
		url:      endpointIdentityLogin,
		body:     &Map{"email": creds.Username, "password": creds.Password, "remember": true},
		response: p,
	}
	err := s.request(req)
	if !p.Authenticated {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	return p, err
}

func (s *Sharesies) reAuthenticate() (*ProfileResponse, error) {
	p := &ProfileResponse{}
	req := &sharesiesRequest{
		method:   http.MethodPost,
		url:      endpointIdentityReAuth,
		body:     &Map{"password": s.creds.Password, "acting_as_id": s.ctx.profile.UserList[0].ID},
		response: p,
	}
	err := s.request(req)
	if err != nil {
		return nil, err
	}

	if !p.Authenticated {
		return nil, fmt.Errorf("failed to re-authenticate: %w", err)
	}

	token, _, _ := new(jwt.Parser).ParseUnverified(p.DistillToken, jwt.MapClaims{})
	s.ctx.token = token
	s.ctx.profile = p

	return p, nil
}

func (s *Sharesies) request(request *sharesiesRequest) error {
	b := &bytes.Buffer{}
	e := json.NewEncoder(b)
	e.Encode(request.body)

	req, err := http.NewRequest(request.method, request.url, b)
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

	res, err := s.httpClient.Do(req)
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
