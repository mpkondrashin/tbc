package sms

import "encoding/xml"

type GetFilters struct {
	XMLName xml.Name `xml:"getFiltes"`
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
