package sms

import (
	"fmt"
	"io"
	"net/http"
)

func (s *SMS) GetActionSetRefID(actionSetName string) (*string, error) {
	client := s.getClient()
	url := s.url + "/dbAccess/tptDBServlet?method=DataDictionary&table=ACTIONSET&format=xml"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s: %w", url, ErrByCode(resp.StatusCode))
	}
	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	d := string(xmlData)
	return &d, nil
}
