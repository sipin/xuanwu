package models

import (
	//Official libs
	"errors"
	"fmt"
	"strconv"
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

func WidgetDisable(ws ...*Widget) {
	for _, w := range ws {
		w.Disable()
	}
}

func WidgetHidden(ws ...*Widget) {
	for _, w := range ws {
		w.Hide()
	}
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
	Required    bool
	Disabled    bool
	Hidden      bool
	GetBindData func() (data []*IDLabelPair)
}

func (w *Widget) Disable() {
	w.Disabled = true
}

func (w *Widget) Hide() {
	w.Hidden = true
}

func (w *Widget) String() string {

	if idx, err := strconv.Atoi(w.Value); err == nil {
		if w.StringList != nil {
			return w.StringList[idx]
		}
		if w.EnumData != nil {
			return w.EnumData[int32(idx)]
		}
	}

	if w.GetBindData != nil {
		datas := w.GetBindData()
		for _, data := range datas {
			if w.Value == data.ID {
				return data.Label
			}
		}
	}

	return w.Value
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
var DateLayout = "2006-01-02"

func I64Time(c int64) string {
	if c == 0 {
		return ""
	}
	return time.Unix(c, 0).Format(DateTimeLayout)
}

func I64Date(c int64) string {
	if c == 0 {
		return ""
	}
	return time.Unix(c, 0).Format(DateLayout)
}

func I32Time(c int32) string {
	return I64Time(int64(c))
}
