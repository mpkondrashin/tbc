package sms

import "encoding/xml"

// Request

type GetFilters struct {
	XMLName xml.Name `xml:"getFilters"`
	Profile Profile  `xml:"profile"`
	Filter  []Filter `xml:"filter"`
}

type SetFilters struct {
	XMLName xml.Name `xml:"setFilters"`
	Profile Profile  `xml:"profile,omitempty"`
	Filter  []Filter `xml:"filter"`
}

type Profile struct {
	ID   string `xml:"id,attr,omitempty"`
	Name string `xml:"name,attr,omitempty"`
}

type PolicyGroup struct {
	//Text  string `xml:",chardata"`
	Refid string `xml:"refid,attr"`
}

type Trigger struct {
	//Text    string `xml:",chardata"`
	Timeout string `xml:"timeout,attr"`
}

type Capability struct {
	//	Text      string `xml:",chardata"`
	Name      string `xml:"name,attr"`
	Enabled   string `xml:"enabled,omitempty"`
	Actionset string `xml:"actionset,omitempty"`
}

// 052 457 1815

type Filter struct {
	Name        string       `xml:"name,omitempty"`
	Number      string       `xml:"number,omitempty"`
	SignatureID string       `xml:"signature-id,omitempty"`
	PolicyID    string       `xml:"policy-id,omitempty"`
	Version     string       `xml:"version,omitempty"`
	Locked      string       `xml:"locked,omitempty"`
	UseParent   string       `xml:"useParent,omitempty"`
	Comment     string       `xml:"comment,omitempty"`
	Description string       `xml:"description,omitempty"`
	Severity    string       `xml:"severity,omitempty"`
	Enabled     string       `xml:"enabled,omitempty"`
	Actionset   string       `xml:"actionset,omitempty"`
	Control     string       `xml:"control,omitempty"`
	Afc         string       `xml:"afc,omitempty"`
	PolicyGroup PolicyGroup  `xml:"policyGroup,omitempty"`
	Trigger     Trigger      `xml:"trigger"`
	Capability  []Capability `xml:"capability"`
}

type Filters struct {
	XMLName                   xml.Name `xml:"filters"`
	Text                      string   `xml:",chardata"`
	NoNamespaceSchemaLocation string   `xml:"noNamespaceSchemaLocation,attr"`
	Xsi                       string   `xml:"xsi,attr"`
	Profile                   struct {
		Text    string `xml:",chardata"`
		Name    string `xml:"name,attr"`
		ID      string `xml:"id,attr"`
		Version string `xml:"version,attr"`
	} `xml:"profile"`
	Filter []Filter `xml:"filter,omitempty"`
}
