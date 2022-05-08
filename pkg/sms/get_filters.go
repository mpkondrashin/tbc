package sms

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

func (s *SMS) GetFilters(getFilters *GetFilters) (*Filters, error) {
	client := s.getClient()
	url := s.url + "/ipsProfileMgmt/getFilters"
	fmt.Println("URL:", url)
	bodyXML, err := xml.Marshal(getFilters)
	if err != nil {
		return nil, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	partHeaders := textproto.MIMEHeader{}
	partHeaders.Set("Content-Type", "application/xml")
	w, err := writer.CreateFormFile("name", "getFilter.xml")
	if err != nil {
		return nil, err
	}

	_, _ = w.Write(bodyXML)
	_ = writer.Close()
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	s.auth.Auth(req)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("User-Agent", s.userAgent)

	/*
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Dump: %s\n\n", string(dump))
	*/
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
	var result Filters
	err = xml.Unmarshal(xmlData, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
