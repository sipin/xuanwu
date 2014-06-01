package models

import (
	//Official libs
	"fmt"
	"errors"

	//3rd party libs
	"github.com/sipin/gothrift/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf

var GoUnusedProtection__ int
var ErrInvalidObjectId = errors.New("Invalid input to ObjectIdHex")

type Widget struct {
	Label       string
	Value       string
	Name        string
	PlaceHolder string
	Type        string
	ErrorMsg    string
}
