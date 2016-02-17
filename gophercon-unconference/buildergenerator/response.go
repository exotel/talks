package exoml

import (
	"encoding/xml"
	"fmt"
)

//Response struct for the verb Response
//the wrapper struct that accomodates all the Verbs as a sequence of verbs
type Response struct {
	XMLName  xml.Name `xml:"Response"`
	Response []interface{}
}

//NewResponse returns a pointer to the reposne structure
func NewResponse() *Response {
	return new(Response)
}

// Action appends action verb structs to response.
//The verbs has to be given in the order in which
//they are expected to be executed
// if there is any invalid verb the function will repond with error but still Response would have all the other
//verbs till the invalid one
func (r *Response) Action(verbs ...interface{}) error {
	for _, s := range verbs {
		r.Response = append(r.Response, s)
	}
	return nil
}

// String returns a formatted xml response
// String implements the fmt.Stringer and it returns the ExoML Response struct as an XMLMarshalled string
func (r Response) String() string {
	output, err := xml.MarshalIndent(r, "", "   ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return ""
	}
	return xml.Header + string(output) + "\n"
}
