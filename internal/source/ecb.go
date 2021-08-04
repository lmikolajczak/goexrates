package source

import "encoding/xml"

type Rate struct {
	Iso   string  `xml:"currency,attr"`
	Value float32 `xml:"rate,attr"`
}

type Day struct {
	Date  string `xml:"time,attr"`
	Rates []Rate `xml:"Cube"`
}

type ECB struct {
	XMLName xml.Name
	Days    []Day `xml:"Cube>Cube"`
}
