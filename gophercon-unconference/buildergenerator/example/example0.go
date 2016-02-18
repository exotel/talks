// +build OMIT

package main

import (
	"fmt"

	"github.com/exotel/go-talks/gophercon-unconference/buildergenerator"
)

func main() {
	dial := exoml.NewDial()
	dial.SetTimeout(5)
	fmt.Println(dial.GetTimeout())
}
