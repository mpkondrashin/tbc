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

type SetFilters struct {
	XMLName xml.Name `xml:"setFilters"`
	Profile struct {
		Name string `xml:"name,attr"`
		ID   string `xml:"id,attr"`
	} `xml:"profile"`
	Filter []struct {
		Number    string `xml:"number,omitempty"`
		Locked    string `xml:"locked,omitempty"`
		Comment   string `xml:"comment,omitempty"`
		Control   string `xml:"control,omitempty"`
		Actionset struct {
			Name  string `xml:"name,attr"`
			Refid string `xml:"refid,attr"`
		} `xml:"actionset"`
		Enabled   string `xml:"enabled,omitempty"`
		Afc       string `xml:"afc,omitempty"`
		UseParent string `xml:"useParent,omitempty"`
		Trigger   struct {
			Threshold string `xml:"threshold,attr"`
			Timeout   string `xml:"timeout,attr"`
		} `xml:"trigger"`
		Name string `xml:"name,omitempty"`
	} `xml:"filter"`
}
