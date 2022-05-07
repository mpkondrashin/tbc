package sms

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

func (s *SMS) GetActionSet(actionSetName string) (*Resultset, error) {
	actionSetName = strings.ReplaceAll(actionSetName, "/", "+")
	fmt.Println("actionSetName", actionSetName)
	client := s.getClient()
	url := s.url + "/dbAccess/tptDBServlet?method=DataDictionary&table=ACTIONSET&format=xml"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	s.auth.Auth(req)
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Dump: %s\n\n", string(dump))

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
	resultset, err := s.GetActionSet(actionSetName)
	if err != nil {
		return "", err
	}
	for _, r := range resultset.Table.Data.R {
		fmt.Printf("compare \"%s\" to \"%s\"\n", r.C[1], actionSetName)
		if r.C[1] == actionSetName {
			return r.C[0], nil
		}
	}
	return "", fmt.Errorf("actionSet \"%s\" not found", actionSetName)
}

type Resultset struct {
	XMLName xml.Name `xml:"resultset"`
	Text    string   `xml:",chardata"`
	Table   struct {
		Text   string `xml:",chardata"`
		Name   string `xml:"name,attr"`
		Column []struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
			Type string `xml:"type,attr"`
		} `xml:"column"`
		Data struct {
			Text string `xml:",chardata"`
			R    []struct {
				Text string   `xml:",chardata"`
				C    []string `xml:"c"`
			} `xml:"r"`
		} `xml:"data"`
	} `xml:"table"`
}
