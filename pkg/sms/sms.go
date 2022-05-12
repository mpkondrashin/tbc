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

func (s *SMS) GetFilters(getFilters *GetFilters) (result *Filters, err error) {
	err = s.SendRequest("POST", "/ipsProfileMgmt/getFilters", getFilters, &result)
	return
}

func (s *SMS) SetFilters(setFilters *SetFilters) error {
	return s.SendRequest("POST", "/ipsProfileMgmt/setFilters", setFilters, nil)
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
		//fmt.Printf("compare \"%s\" to \"%s\"\n", r.C[2], action)
		log.Printf("ActionSet: %s %s %s %s", r.C[0], r.C[1], r.C[2], r.C[3])
		if r.C[2] == action {
			result = append(result, r.C[0])
			log.Printf("ActionSet with %s action: %s", action, r.C[1])
		}
	}
	return result, nil
}

func (s *SMS) DistributeProfileX(distribution *Distribution) error {
	return s.SendRequest("POST", "/ipsProfileMgmt/distributeProfile", distribution, nil)
}

func (s *SMS) DataDictionaryX(table string) (result *Resultset, err error) {
	url := fmt.Sprintf("/dbAccess/tptDBServlet?method=DataDictionary&table=%s&format=xml", table)
	err = s.SendRequest("GET", url, nil, result)
	return
}

func (s *SMS) DataDictionaryAll() (err error) {
	var result interface{}
	err = s.SendRequest("GET", "/dbAccess/tptDBServlet?method=DataDictionary&format=xml", nil, &result)
	fmt.Println("ERR ", err)
	fmt.Println("RESULT", result)
	return nil
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

/*
func (s *SMS) GetCategory() (*Resultset, error) {
	return s.DataDictionary("CATEGORY")
}
*/

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
	status := SeekForStatus(string(xmlData))
	if status != "" {
		return fmt.Errorf("SMS: %s", status)
	}
	if reply != nil {
		return xml.Unmarshal(xmlData, reply)
	}
	return nil
}

func SeekForStatus(data string) string {
	start := strings.Index(data, "<status>")
	if start == -1 {
		return ""
	}
	end := strings.Index(data, "</status>")
	if end == -1 {
		return ""
	}
	if end < start {
		return ""
	}
	return data[start+8 : end]
}

func (s *SMS) GetFilters_(getFilters *GetFilters) (*Filters, error) {
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

	fmt.Println("Get Filters: ", string(xmlData))
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

func (s *SMS) SetFilters_(setFilters *SetFilters) error {
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

func (s *SMS) GetDistributionStatus(distribution *Distribution) error {
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
