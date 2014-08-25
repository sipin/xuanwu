package xuanwu

import (
	"strconv"
	"strings"
	"time"
)

type IDLabelPair struct {
	ID    string
	Label string
	Order string
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
	Label          string
	Value          string
	Name           string
	PlaceHolder    string
	Type           string
	ErrorMsg       string
	EnumKey        []int32
	EnumData       map[int32]string
	StringList     []string
	Required       bool
	Disabled       bool
	Hidden         bool
	Readonly       bool
	IsList         bool
	Option         map[string]bool
	GetBindData    func() (data []*IDLabelPair)
	GetMetaData    func() (data []IXuanWuObj, head IXuanWuObj)
	GetGroupedData func() (subInput string, data map[string][]*IDLabelPair)
}

func (w *Widget) SetOption(name string) {
	if w.Option == nil {
		w.Option = make(map[string]bool, 2)
	}
	w.Option[name] = true
}

func (w *Widget) HasOption(name string) bool {
	return w.Option[name]
}

func (w *Widget) Disable() {
	w.Disabled = true
}

func (w *Widget) Hide() {
	w.Hidden = true
}

func (w *Widget) Val() string {

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
		if w.IsList {
			result := []string{}
			values := strings.Split(w.Value, "\n")
			for _, value := range values {
				for _, data := range datas {
					if value == data.ID {
						result = append(result, data.Label)
						break
					}
				}
			}
			return strings.Join(result, "\n")
		} else {
			for _, data := range datas {
				if w.Value == data.ID {
					return data.Label
				}
			}
		}
	}

	if w.GetMetaData != nil {
		// should show not thing
		return ""
	}

	return w.Value
}

type IXuanWuObj interface {
	Id() string
	IsSearchEnabled() bool
	GetLabel() string
	GetName() string
	GetListedLabels() []*IDLabelPair
	GetFieldAsString(fieldKey string) (Value string)
	GetFilters() []*Widget
	Widgets() []*Widget
}

var DateTimeLayout = "2006-01-02 15:04"
var DateLayout = "2006-01-02"
var TimeLayout = "15:04"

func I64DateTime(c int64) string {
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

func I64Time(c int64) string {
	if c == 0 {
		return ""
	}
	return time.Unix(c, 0).Format(TimeLayout)
}

func I32Time(c int32) string {
	return I64Time(int64(c))
}

func XGetQueryString(word string, fields []string) map[string]interface{} {
	queryString := make(map[string]interface{})
	queryString["default_operator"] = "AND"
	queryString["fields"] = fields
	queryString["query"] = word

	return queryString
}

func XGetQuery(key string, data map[string]interface{}) map[string]interface{} {
	query := make(map[string]interface{})
	query[key] = data
	args := make(map[string]interface{})
	args["query"] = query
	return args
}

func parseTime(layout, str string) time.Time {
	now := time.Now()
	t, _ := time.ParseInLocation(layout, str, now.Location())
	return t
}

func XGetSearchObj(word string, fields []string, params map[string]string, termKeys map[string]bool, dateKeys map[string]bool) map[string]interface{} {
	terms := make(map[string]string)
	ranges := make(map[string]map[string]int64)

	for k, v := range params {
		if v == "" {
			continue
		}
		if _, ok := termKeys[k]; ok {
			terms[k] = v
			continue
		}

		if isStart, ok := dateKeys[k]; ok {
			intVal := parseTime(DateLayout, v)
			if isStart {
				fieldName := k[0 : len(k)-5]
				if dateVal, ok := ranges[fieldName]; ok {
					dateVal["gte"] = intVal.Unix()
					ranges[fieldName] = dateVal
				} else {
					ranges[fieldName] = map[string]int64{
						"gte": intVal.Unix(),
						"lt":  intVal.AddDate(0, 0, 1).Unix(),
					}
				}
			} else {
				fieldName := k[0 : len(k)-3]
				if dateVal, ok := ranges[fieldName]; ok {
					dateVal["lt"] = intVal.AddDate(0, 0, 1).Unix()
				} else {
					ranges[fieldName] = map[string]int64{
						"gte": intVal.AddDate(0, 0, -1).Unix() + 1,
						"lt":  intVal.AddDate(0, 0, 1).Unix(),
					}
				}
			}
		}
	}

	if len(terms) == 0 && len(ranges) == 0 {
		return XGetQuery("query_string", XGetQueryString(word, fields))
	}

	filtered := make(map[string]interface{})
	if word != "" {
		filtered["query"] = map[string]interface{}{
			"query_string": XGetQueryString(word, fields),
		}
	}

	filter := make(map[string]interface{})
	var must []interface{}

	for k, v := range terms {
		must = append(must, map[string]interface{}{
			"term": map[string]string{
				k: v,
			},
		})
	}

	for k, v := range ranges {
		must = append(must, map[string]interface{}{
			"range": map[string]interface{}{
				k: v,
			},
		})
	}
	filter["bool"] = map[string]interface{}{
		"must": must,
	}
	filtered["filter"] = filter

	return XGetQuery("filtered", filtered)
}

func XGetMoreSearchObj(word string, fields []string, params map[string]interface{}, termKeys map[string]bool, dateKeys map[string]bool) map[string]interface{} {
	terms := make(map[string]interface{})
	ranges := make(map[string]map[string]int64)

	for k, v := range params {
		switch v := v.(type) {
		case string:
			if _, ok := termKeys[k]; ok {
				terms[k] = v
				continue
			}

			if isStart, ok := dateKeys[k]; ok {
				intVal := parseTime(DateLayout, v)
				if isStart {
					fieldName := k[0 : len(k)-5]
					if dateVal, ok := ranges[fieldName]; ok {
						dateVal["gte"] = intVal.Unix()
						ranges[fieldName] = dateVal
					} else {
						ranges[fieldName] = map[string]int64{
							"gte": intVal.Unix(),
							"lt":  intVal.AddDate(0, 0, 1).Unix(),
						}
					}
				} else {
					fieldName := k[0 : len(k)-3]
					if dateVal, ok := ranges[fieldName]; ok {
						dateVal["lt"] = intVal.AddDate(0, 0, 1).Unix()
					} else {
						ranges[fieldName] = map[string]int64{
							"gte": intVal.AddDate(0, 0, -1).Unix() + 1,
							"lt":  intVal.AddDate(0, 0, 1).Unix(),
						}
					}
				}
			}
		case []string:
			if len(v) == 0 {
				continue
			}
			if _, ok := termKeys[k]; ok {
				terms[k] = v
				continue
			}
		}
	}

	if len(terms) == 0 && len(ranges) == 0 {
		return XGetQuery("query_string", XGetQueryString(word, fields))
	}

	filtered := make(map[string]interface{})
	if word != "" {
		filtered["query"] = map[string]interface{}{
			"query_string": XGetQueryString(word, fields),
		}
	}

	filter := make(map[string]interface{})
	var must []interface{}
	var should []interface{}

	for k, v := range terms {
		switch v := v.(type) {
		case string:
			must = append(must, map[string]interface{}{
				"term": map[string]string{
					k: v,
				},
			})
		case []string:
			for _, val := range v {
				should = append(should, map[string]interface{}{
					"term": map[string]string{
						k: val,
					},
				})
			}
		}
	}

	for k, v := range ranges {
		must = append(must, map[string]interface{}{
			"range": map[string]interface{}{
				k: v,
			},
		})
	}
	filter["bool"] = map[string]interface{}{
		"must":   must,
		"should": should,
	}
	filtered["filter"] = filter

	return XGetQuery("filtered", filtered)
}

func XSortFieldsFilter(sortFields []string) (rtn []string) {
	rtn = make([]string, 0, len(sortFields))
	for _, s := range sortFields {
		if s != "" {
			rtn = append(rtn, s)
		}
	}
	return
}
