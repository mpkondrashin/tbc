package sms

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (s *SMS) GetFilters(getFilters *GetFilters) (*string, error) {
	client := s.getClient()
	url := s.url + "/ipsProfileMgmt/getFilters"
	bodyXML, err := xml.Marshal(getFilters)

	if err != nil {
		return nil, err
	}

	fmt.Println(string(bodyXML))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyXML))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	s.auth.Auth(req)
	//req.Header.Add("Accept", "application/xml")
	req.Header.Add("User-Agent", s.userAgent)

	fmt.Println(formatRequest(req))
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

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	var request []string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	return strings.Join(request, "\n")
}
