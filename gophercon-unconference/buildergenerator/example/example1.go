// +build OMIT

package main

import (
	"fmt"

	"github.com/exotel/talks/gophercon-unconference/buildergenerator"
)

func main() {
	resp := exoml.NewResponse()

	dial := exoml.NewDial().SetTimeout(5)

	_ = resp.Action(dial)
	fmt.Println(resp)
}
