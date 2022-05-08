package sms

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

func (s *SMS) SetFilters(setFilters *SetFilters) error {
	bodyXML, err := xml.Marshal(setFilters)
	if err != nil {
		return err
	}
	client := s.getClient()
	url := s.url + "/ipsProfileMgmt/setFilters"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	partHeaders := textproto.MIMEHeader{}
	partHeaders.Set("Content-Type", "application/xml")
	w, err := writer.CreateFormFile("name", "setFilter.xml")
	if err != nil {
		return err
	}
	_, _ = w.Write(bodyXML)
	_ = writer.Close()
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}
	s.auth.Auth(req)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("User-Agent", s.userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http.Client.Do: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("%s: %w", url, ErrByCode(resp.StatusCode))
	}
	return nil
}
