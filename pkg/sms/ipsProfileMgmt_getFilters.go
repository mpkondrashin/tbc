package sms

import "encoding/xml"

type GetFilters struct {
	XMLName xml.Name `xml:"getFilters"`
	Profile Profile  `xml:"profile"`
	Filter  []Filter `xml:"filter"`
}

type Profile struct {
	ID   string `xml:"id,attr,omitempty"`
	Name string `xml:"name,attr,omitempty"`
}

type Filter struct {
	Number      uint   `xml:"number,omitempty"`
	Name        string `xml:"name,omitempty"`
	SignatureID string `xml:"signature-id,omitempty"`
	PolicyID    string `xml:"policy-id,omitempty"`
}

// Response

type Filters struct {
	XMLName xml.Name `xml:"filters"`
	Text    string   `xml:",chardata"`
	Profile struct {
		Text    string `xml:",chardata"`
		Name    string `xml:"name,attr"`
		ID      string `xml:"id,attr"`
		Version string `xml:"version,attr"`
	} `xml:"profile"`
	Filter struct {
		Text      string `xml:",chardata"`
		Name      string `xml:"name"`
		PolicyID  string `xml:"policy-id"`
		Version   string `xml:"version"`
		Locked    string `xml:"locked"`
		UseParent string `xml:"useParent"`
		Comment   string `xml:"comment"`
		Severity  string `xml:"severity"`
		Enabled   string `xml:"enabled"`
		Actionset struct {
			Text  string `xml:",chardata"`
			Refid string `xml:"refid,attr"`
			Name  string `xml:"name,attr"`
		} `xml:"actionset"`
		Control string `xml:"control"`
		Afc     string `xml:"afc"`
	} `xml:"filter"`
}
