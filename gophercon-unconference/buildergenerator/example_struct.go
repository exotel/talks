package exoml

import "encoding/xml"

//go:generate generator -t Dial -f $GOFILE

type Dial struct {
	XMLName      xml.Name `xml:"Dial"`
	Action       string   `xml:"action,attr"`
	Method       string   `xml:"method,attr"`
	Timeout      int      `xml:"timeout,attr"`
	HangupOnStar bool     `xml:"hangupOnStar,attr"`
	TimeLimit    int      `xml:"timeLimit,attr"`
	CallerID     string   `xml:"callerId,attr"`
	Record       bool     `xml:"record,attr"`
}
