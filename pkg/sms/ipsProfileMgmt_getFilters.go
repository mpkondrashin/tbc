package sms

type getFilters struct {
	Profile profile  `xml:"profile"`
	Filter  []filter `xml:"filter"`
}

type profile struct {
	ID   string `xml:"id,attr,omitempty"`
	Name string `xml:"name,attr,omitempty"`
}

type filter struct {
	Number      uint   `xml:"number,omitempty"`
	Name        string `xml:"name,omitempty"`
	SignatureID string `xml:"signature-id,omitempty"`
	PolicyID    string `xml:"policy-id,omitempty"`
}
