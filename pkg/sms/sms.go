package sms

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

/*

var (
	//	ErrTooBigFile    = errors.New("too big file size")
	ErrResponseError = errors.New("response error")
)
*/
const version = "0.1"
const defaultUserAgent = "smsClient/" + version

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

func (s *SMS) GetFilters(getFilters *GetFilters) (*Filters, error) {
	client := s.getClient()
	url := s.url + "/ipsProfileMgmt/getFilters"
	//ntln("URL:", url)
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
	//fmt.Println("Response:", resp)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s: %w", url, ErrByCode(resp.StatusCode))
	}
	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	//fmt.Println("Get Filters: ", string(xmlData))
	var result Filters
	err = xml.Unmarshal(xmlData, &result)
	if err != nil {
		return nil, err
	}
	if result.Status != nil {
		return nil, fmt.Errorf("GetFilters: %s", result.Status.Text)
	}
	if len(result.Filter) > 0 && result.Filter[0].Status != nil {
		return nil, fmt.Errorf("GetFilters: %s", result.Filter[0].Status.Text)
	}
	return &result, nil
}

func (s *SMS) SetFilters(setFilters *SetFilters) error {
	bodyXML, err := xml.Marshal(setFilters)
	if err != nil {
		return err
	}
	//fmt.Println("bodyXML", string(bodyXML))
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

func (s *SMS) GetActionSet() (*Resultset, error) {
	client := s.getClient()
	url := s.url + "/dbAccess/tptDBServlet?method=DataDictionary&table=ACTIONSET&format=xml"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	s.auth.Auth(req)
	/*
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			panic(err)
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

	var result Resultset
	err = xml.Unmarshal(xmlData, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SMS) GetActionSetRefID(actionSetName string) (string, error) {
	actionSetName = strings.ReplaceAll(actionSetName, "/", "+")
	resultset, err := s.GetActionSet()
	if err != nil {
		return "", err
	}
	for _, r := range resultset.Table.Data.R {
		//fmt.Printf("compare \"%s\" to \"%s\"\n", r.C[1], actionSetName)
		if r.C[1] == actionSetName {
			return r.C[0], nil
		}
	}
	return "", fmt.Errorf("actionSet \"%s\" not found", actionSetName)
}

type DistributionPiority int

const (
	PriorityLow DistributionPiority = iota
	PriorityHigh
)

var DistributionPiorityString = []string{"low", "high"}

func (d DistributionPiority) String() string {
	return DistributionPiorityString[d]
}

func DistributionPiorityFromString(s string) (DistributionPiority, error) {
	switch s {
	case "low":
		return PriorityLow, nil
	case "high":
		return PriorityHigh, nil
	default:
		return PriorityLow, fmt.Errorf("Unknown priority \"%s\"", s)
	}
}

func (s *SMS) DistributeProfile(distribution *Distribution) error {
	bodyXML, err := xml.Marshal(distribution)
	if err != nil {
		return err
	}
	fmt.Println("DistributeProfile, bodyXML", string(bodyXML))
	client := s.getClient()
	url := s.url + "/ipsProfileMgmt/distributeProfile"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	partHeaders := textproto.MIMEHeader{}
	partHeaders.Set("Content-Type", "application/xml")
	w, err := writer.CreateFormFile("name", "distributeProfile.xml")
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
	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err //nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	fmt.Println("Distribution status:", string(xmlData))
	var result Distribution
	err = xml.Unmarshal(xmlData, &result)
	if err != nil {
		return err
	}
	if result.SegmentGroup.Status != nil {
		return fmt.Errorf("DistributeProfile: %s", result.SegmentGroup.Status.Text)
	}
	return nil
}

func (s *SMS) GetDistribionStatus(distribution *Distribution) error {
	bodyXML, err := xml.Marshal(distribution)
	if err != nil {
		return err
	}
	fmt.Println("DistributeProfile, bodyXML", string(bodyXML))
	client := s.getClient()
	url := s.url + "/ipsProfileMgmt/distributionStatus"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	partHeaders := textproto.MIMEHeader{}
	partHeaders.Set("Content-Type", "application/xml")
	w, err := writer.CreateFormFile("name", "distributeProfile.xml")
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
