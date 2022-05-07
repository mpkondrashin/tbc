package sms

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
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

	fmt.Println(string(bodyXML))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	partHeaders := textproto.MIMEHeader{}
	partHeaders.Set("Content-Type", "application/xml")
	//w, err := writer.CreatePart(partHeaders)
	w, err := writer.CreateFormFile("name", "getFilter.xml")
	//w, err := writer.CreateFormField("xml")
	if err != nil {
		return nil, err
	}
	//bodyXML = []byte("<getFilters><profile name=\"tbcheck\"/><filter><number>51</number></filter></getFilters>")
	//body := bytes.Buffer{} bodyXML
	_, _ = w.Write(bodyXML)
	_ = writer.Close()
	req, err := http.NewRequest("POST", url, body)

	//req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyXML))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	s.auth.Auth(req)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", writer.FormDataContentType())
	//req.Header.Add("Content-Type", "application/xml")
	req.Header.Add("User-Agent", s.userAgent)

	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Dump: %s\n\n", string(dump))

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
	var result Filters
	err = xml.Unmarshal(xmlData, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
