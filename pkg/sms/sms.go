package sms

import (
	"crypto/tls"
	"net/http"
)

/*

var (
	//	ErrTooBigFile    = errors.New("too big file size")
	ErrResponseError = errors.New("response error")
)
*/
const defaultUserAgent = "github.com/mpkondrashin/tbcheck/pkg/sms"

type SMS struct {
	url                string
	auth               Autherization
	userAgent          string
	insecureSkipVerify bool
}

func New(url string, auth Autherization) *SMS {
	return &SMS{
		url:       url,
		auth:      auth,
		userAgent: defaultUserAgent,
	}
}

func (s *SMS) SetUserAgent(userAgent string) *SMS {
	s.userAgent = userAgent
	return s
}

func (s *SMS) SetInsecureSkipVerify(insecureSkipVerify bool) *SMS {
	s.insecureSkipVerify = insecureSkipVerify
	return s
}

func (s *SMS) getClient() *http.Client {
	if s.insecureSkipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		return &http.Client{Transport: tr}
	}
	return &http.Client{}
}
