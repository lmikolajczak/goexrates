package source

import "encoding/xml"

type Currency struct {
	Code  string  `xml:"currency,attr"`
	Value float64 `xml:"rate,attr"`
}

type Day struct {
	Date       string     `xml:"time,attr"`
	Currencies []Currency `xml:"Cube"`
}

type ECB struct {
	XMLName xml.Name
	Days    []Day `xml:"Cube>Cube"`
}
