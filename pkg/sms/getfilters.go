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
)

func (s *SMS) GetFilters(getFilters *GetFilters) (*string, error) {
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
	w, err := writer.CreateFormFile("name", "getFilter.xml")
	//w, err := writer.CreateFormField("xml")
	if err != nil {
		return nil, err
	}
	w.Write(bodyXML)

	//body := bytes.Buffer{} bodyXML
	req, err := http.NewRequest("POST", url, body)
	//req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("Somecrap")))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	s.auth.Auth(req)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("User-Agent", s.userAgent)

	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", string(dump))

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
