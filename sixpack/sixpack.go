package sixpack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

const pattern string = `^[a-z0-9][a-z0-9\-_ ]*$`

var BaseUrl string = "http://localhost:5000"

type Options struct {
	ClientID  []byte
	BaseUrl   *url.URL
	IP        string
	UserAgent string
}

func (opts *Options) Values() url.Values {
	v := url.Values{}

	if opts.IP != "" {
		v.Set("ip_address", opts.IP)
	}

	if opts.UserAgent != "" {
		v.Set("user_agent", opts.UserAgent)
	}

	return v
}

type Alternative struct {
	Name string `json:"name"`
}

type Experiment struct {
	Version int    `json:"version"`
	Name    string `json:"name"`
}

type Response struct {
	Status      string      `json:"status"`
	ClientID    string      `json:"client_id"`
	Alternative Alternative `json:"alternative"`
	Experiment  Experiment  `json:"experiment"`
}

type Session struct {
	Client  *http.Client
	Options *Options
}

func NewSession(opts Options) (*Session, error) {
	if len(opts.ClientID) == 0 {
		id, err := GenerateClientID()
		if err != nil {
			return nil, err
		}
		opts.ClientID = id
	}

	if opts.BaseUrl == nil {
		u, err := url.Parse(BaseUrl)
		if err != nil {
			return nil, err
		}
		opts.BaseUrl = u
	}

	return &Session{Client: &http.Client{}, Options: &opts}, nil
}

func (s *Session) Participate(name string, alternatives []string, force string) (r *Response, err error) {
	rex := regexp.MustCompile(pattern)
	if !rex.MatchString(name) {
		return nil, errors.New("Bad experiment name")
	}

	if len(alternatives) < 2 {
		return nil, errors.New("Must specify at least 2 alternatives")
	}

	for _, alt := range alternatives {
		if !rex.MatchString(alt) {
			return nil, errors.New(fmt.Sprintf("Bad alternative name: %s", alt))
		}
	}

	if force != "" {
		for _, alt := range alternatives {
			if force == alt {
				return s.forceResponse(name, force), nil
			}
		}
	}

	endpoint, _ := url.Parse("/participate")

	params := s.Options.Values()
	params.Set("client_id", string(s.Options.ClientID))
	params.Set("experiment", name)

	for _, alt := range alternatives {
		params.Add("alternatives", alt)
	}

	r, err = s.request(endpoint, params)
	if err != nil {
		return s.fallbackResponse(alternatives[0]), nil
	}

	return
}

func (s *Session) Convert(name string) (r *Response, err error) {
	endpoint, _ := url.Parse("/convert")

	params := s.Options.Values()
	params.Set("client_id", string(s.Options.ClientID))
	params.Set("experiment", name)

	return s.request(endpoint, params)
}

func (s *Session) request(endpoint *url.URL, params url.Values) (r *Response, err error) {
	u := s.Options.BaseUrl.ResolveReference(endpoint)
	u.RawQuery = params.Encode()

	resp, err := s.Client.Get(u.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return nil, errors.New(string(b))
	}

	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}

	return
}

func (s *Session) forceResponse(name, force string) *Response {
	return &Response{
		Status: "ok",
		Alternative: Alternative{
			Name: force,
		},
		Experiment: Experiment{
			Version: 0,
			Name:    name,
		},
		ClientID: string(s.Options.ClientID),
	}
}

func (s *Session) fallbackResponse(alt string) *Response {
	return &Response{
		Status: "failed",
		Alternative: Alternative{
			Name: alt,
		},
	}
}
