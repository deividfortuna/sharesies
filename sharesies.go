package sharesies

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	"github.com/dgrijalva/jwt-go"
)

const appSharesiesUrl = "https://app.sharesies.nz/api/"
const dataSharesiesUrl = "https://data.sharesies.nz/api/"

type Map map[string]interface{}

type SharesiesCtx struct {
	token   *jwt.Token
	profile *Profile
}

type Sharesies struct {
	httpClient *http.Client
	creds      *SharesiesCredentials
	ctx        *SharesiesCtx
}

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

func NewSharesies(creds *SharesiesCredentials) (*Sharesies, error) {
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
		&SharesiesCtx{},
	}

	errCtx := i.getSharesiesCtx(i.ctx)
	if errCtx != nil {
		return nil, errCtx
	}

	return i, nil
}

func (s *Sharesies) getSharesiesCtx(ctx *SharesiesCtx) error {
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

func (s *Sharesies) Profile() (*Profile, error) {
	r := &Profile{}
	err := s.request(&sharesiesRequest{method: http.MethodGet, url: appSharesiesUrl + "identity/check", response: r})
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *Sharesies) Companies(request *InstrumentsRequest) (*InstrumentResponse, error) {
	r := &InstrumentResponse{}
	h, errH := s.getAuthHeaders()
	if errH != nil {
		return nil, errH
	}

	err := s.request(&sharesiesRequest{method: http.MethodPost, url: dataSharesiesUrl + "v1/instruments", body: request, response: r, headers: h})
	return r, err
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

func (s *Sharesies) authenticate(creds *SharesiesCredentials) (*Profile, error) {
	r := &Profile{}
	err := s.request(&sharesiesRequest{method: http.MethodPost, url: appSharesiesUrl + "identity/login", body: &Map{"email": creds.Username, "password": creds.Password, "remember": true}, response: r})
	if !r.Authenticated {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	return r, err
}

func (s *Sharesies) reAuthenticate() (*Profile, error) {
	p := &Profile{}

	err := s.request(&sharesiesRequest{method: http.MethodPost, url: appSharesiesUrl + "identity/reauthenticate", body: &Map{"password": s.creds.Password, "acting_as_id": s.ctx.profile.UserList[0].ID}, response: p})
	if err != nil {
		return nil, err
	}

	if !p.Authenticated {
		return nil, fmt.Errorf("failed to re-authenticate: %w", err)
	}

	token, _, err := new(jwt.Parser).ParseUnverified(p.DistillToken, jwt.MapClaims{})
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
