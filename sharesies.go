package sharesies

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

const appSharesiesUrl = "https://app.sharesies.nz/api/"

type Map map[string]interface{}

type Sharesies struct {
	httpClient *http.Client
	creds      *SharesiesCredentials
}

type SharesiesCredentials struct {
	Username string
	Password string
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
	}

	_, errAuth := i.Authenticate(creds)
	if errAuth != nil {
		return nil, errAuth
	}

	return i, nil
}

func (s *Sharesies) Authenticate(creds *SharesiesCredentials) (*Profile, error) {
	r := &Profile{}
	err := s.request(http.MethodPost, appSharesiesUrl+"identity/login", &Map{"email": creds.Username, "password": creds.Password, "remember": true}, r)
	if !r.Authenticated {
		return nil, fmt.Errorf("Failed to authenticate")
	}
	return r, err
}

func (s *Sharesies) Profile() (*Profile, error) {
	r := &Profile{}
	err := s.request(http.MethodGet, appSharesiesUrl+"identity/check", nil, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *Sharesies) request(method string, url string, body interface{}, response interface{}) error {
	b := &bytes.Buffer{}
	e := json.NewEncoder(b)
	e.Encode(body)

	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 Firefox/71.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/json")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("expected 200 code")
	}

	bd, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	json.Unmarshal(bd, &response)
	return nil
}
