package models

import (
	//Official libs
	"errors"
	"fmt"
	"time"

	//3rd party libs
	"github.com/sipin/gothrift/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf

var GoUnusedProtection__ int
var ErrInvalidObjectId = errors.New("Invalid input to ObjectIdHex")

type IDLabelPair struct {
	ID    string
	Label string
}

type Widget struct {
	Label       string
	Value       string
	Name        string
	PlaceHolder string
	Type        string
	ErrorMsg    string
	EnumData    map[int32]string
	StringList  []string
	GetBindData func() (data []*IDLabelPair)
}

type IXuanWuObj interface {
	Id() string
	HasSimpleSearch() bool
	GetLabel() string
	GetListedLabels() []*IDLabelPair
	GetFieldAsString(fieldKey string) (Value string)
	Widgets() []*Widget
}

var DateTimeLayout = "2006-01-02 15:04"

func I64Time(c int64) string {
	if c == 0 {
		return ""
	}
	return time.Unix(c, 0).Format(DateTimeLayout)
}

func I32Time(c int32) string {
	return I64Time(int64(c))
}
