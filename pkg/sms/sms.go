package sms

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"log"
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

func (s *SMS) SendRequest(method, url string, request, reply interface{}) error {
	client := s.getClient()
	bodyXML, err := xml.Marshal(request)
	if err != nil {
		return err
	}
	body := &bytes.Buffer{}
	contentType := "application/xml"
	if request != nil {
		writer := multipart.NewWriter(body)
		contentType = writer.FormDataContentType()
		partHeaders := textproto.MIMEHeader{}
		partHeaders.Set("Content-Type", "application/xml")
		w, err := writer.CreateFormFile("name", "get.xml")
		if err != nil {
			return err
		}
		_, _ = w.Write(bodyXML)
		_ = writer.Close()
	}
	req, err := http.NewRequest(method, s.url+url, body)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}
	s.auth.Auth(req)
	req.Header.Add("Accept", "*/*")
	if request != nil {
		req.Header.Add("Content-Type", contentType)
	}
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
		return fmt.Errorf("http.Client.Do: %w", err)
	}
	defer resp.Body.Close()
	//fmt.Println("Response:", resp)
	if resp.StatusCode != 200 {
		return fmt.Errorf("%s: %w", url, ErrByCode(resp.StatusCode))
	}
	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}
	var status Status
	err = xml.Unmarshal(xmlData, &status)
	if err == nil {
		if status.Text != "" {
			return fmt.Errorf("Reply: %s", status.Text)
		}
	}
	if reply != nil {
		err = xml.Unmarshal(xmlData, &reply)
		if err != nil {
			return err
		}
	}
	return nil
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
	return s.DataDictionary("ACTIONSET")
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

func (s *SMS) GetActionSetRefIDsForAction(action string) ([]string, error) {
	switch action {
	case "ALLOW", "DENY", "TRUST", "RATE":
	default:
		return nil, fmt.Errorf("unknown action \"%s\"", action)
	}
	resultset, err := s.GetActionSet()
	if err != nil {
		return nil, err
	}
	var result []string
	for _, r := range resultset.Table.Data.R {
		fmt.Printf("compare \"%s\" to \"%s\"\n", r.C[2], action)
		if r.C[2] == action {
			result = append(result, r.C[0])
			log.Printf("ActionSet with %s action: %s", action, r.C[1])
		}
	}
	return result, nil
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
	log.Println("DistributeProfile() Requst:", string(bodyXML))
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

	log.Println("DistributeProfile() Reply:", string(xmlData))
	var result Distributions
	err = xml.Unmarshal(xmlData, &result)
	if err != nil {
		return err
	}
	if result.SegmentGroup != nil && result.SegmentGroup.Status != nil {
		return fmt.Errorf("DistributeProfile: %s", result.SegmentGroup.Status.Text)
	}
	return nil
}

func (s *SMS) GetDistribionStatus(distribution *Distribution) error {
	bodyXML, err := xml.Marshal(distribution)
	if err != nil {
		return err
	}
	log.Println("GetDistribionStatus() Request:", string(bodyXML))
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
	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err //nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	log.Println("GetDistribionStatus() Reply:", string(xmlData))
	return nil
}

func (s *SMS) DataDictionary(table string) (*Resultset, error) {
	client := s.getClient()
	url := fmt.Sprintf("%s/dbAccess/tptDBServlet?method=DataDictionary&table=%s&format=xml", s.url, table)
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

func (s *SMS) GetSegmentGroups() (*Resultset, error) {
	return s.DataDictionary("SEGMENT_GROUP")
}

func (s *SMS) GetSegmentGroupId(name string) (string, error) {
	resultset, err := s.GetSegmentGroups()
	if err != nil {
		return "", err
	}
	for _, r := range resultset.Table.Data.R {
		//fmt.Printf("compare \"%s\" to \"%s\"\n", r.C[1], name)
		if r.C[1] == name {
			return r.C[0], nil
		}
	}
	return "", fmt.Errorf("SegmentGroup \"%s\" not found", name)
}
