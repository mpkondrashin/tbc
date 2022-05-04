package sms

/*

var (
	//	ErrTooBigFile    = errors.New("too big file size")
	ErrResponseError = errors.New("response error")
)
*/
const defaultUserAgent = "github.com/mpkondrashin/tbcheck/pkg/sms"

type SMS struct {
	url       string
	auth      Autherization
	userAgent string
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
