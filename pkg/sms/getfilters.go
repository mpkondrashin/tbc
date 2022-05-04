package sms

import (
	"fmt"
	"io"
	"net/http"
)

func (s *SMS) GetFilters() (*string, error) {
	client := s.getClient()
	// url := s.url + "/ipsProfileMgmt/getFilters"
	url := fmt.Sprintf("%s%s?profile=%s",
		s.url,
		"/ipsProfileMgmt/getFilters",
		"test",
		//"0051",
	)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	s.auth.Auth(req)
	//req.Header.Add("Accept", "application/xml")
	req.Header.Add("User-Agent", s.userAgent)
	//req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.Client.Do: %w", err)
	}
	defer resp.Body.Close()
	//fmt.Printf("Respond: %v", resp)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s: %w", url, ErrByCode(resp.StatusCode))
	}
	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	//fmt.Printf("%v\n", string(jsonData))
	//var data ListLatest
	//err = json.Unmarshal(jsonData, &data)
	//if err != nil {
	//	return nil, fmt.Errorf("json.Unmarshal: %w: %v\n%s", ErrResponseError, err, string(jsonData))
	//}
	xmlDataStr := string(xmlData)
	return &xmlDataStr, nil
}
