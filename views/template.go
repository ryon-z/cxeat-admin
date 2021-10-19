package views

import (
	"html/template"
	"reflect"
	"strconv"
	"strings"
	"time"
	"yelloment-admin/models"
)

// FuncMap 템플릿 함수맵
func FuncMap() template.FuncMap {
	gfm := make(map[string]interface{}, len(genericMap))
	for k, v := range genericMap {
		gfm[k] = v
	}
	return template.FuncMap(gfm)
}

var genericMap = map[string]interface{}{
	"hello": func() string { return "Hello!" },

	"codelabel":    codeLabel,
	"sepcodelabel": sepCodeLabel,

	"default": dfault,
	"empty":   empty,

	// Date functions
	"date": date,
	"now":  time.Now,

	// Strings
	"trim":  strings.TrimSpace,
	"upper": strings.ToUpper,
	"lower": strings.ToLower,
	"title": strings.Title,
	// Switch order so that "foo" | repeat 5
	"repeat": func(count int, str string) string { return strings.Repeat(str, count) },
	// Deprecated: Use trimAll.
	"trimall": func(a, b string) string { return strings.Trim(b, a) },
	// Switch order so that "$foo" | trimall "$"
	"trimAll":    func(a, b string) string { return strings.Trim(b, a) },
	"trimSuffix": func(a, b string) string { return strings.TrimSuffix(b, a) },
	"trimPrefix": func(a, b string) string { return strings.TrimPrefix(b, a) },
	// Switch order so that "foobar" | contains "foo"
	"contains":    func(substr string, str string) bool { return strings.Contains(str, substr) },
	"containsarr": func(substr string, str []string) bool { return arrContains(str, substr) },
	"hasPrefix":   func(substr string, str string) bool { return strings.HasPrefix(str, substr) },
	"hasSuffix":   func(substr string, str string) bool { return strings.HasSuffix(str, substr) },

	"atoi": func(a string) int { i, _ := strconv.Atoi(a); return i },
}

func arrContains(arrstring []string, substr string) bool {
	for _, s := range arrstring {
		if s == substr {
			return true
		}
	}
	return false
}

func codeLabel(codeType string, codekey string) string {
	return models.GetCodeLabel(codeType, codekey)
}

func sepCodeLabel(codeType string, codekeys string) string {
	return models.GetCodeLabelSeparate(codeType, codekeys)
}

func date(fmt string, date interface{}) string {
	return dateInZone(fmt, date, "Local")
}

func dateInZone(fmt string, date interface{}, zone string) string {
	var t time.Time
	switch date := date.(type) {
	default:
		t = time.Now()
	case time.Time:
		t = date
	case *time.Time:
		t = *date
	case int64:
		t = time.Unix(date, 0)
	case int:
		t = time.Unix(int64(date), 0)
	case int32:
		t = time.Unix(int64(date), 0)
	}

	loc, err := time.LoadLocation(zone)
	if err != nil {
		loc, _ = time.LoadLocation("UTC")
	}

	return t.In(loc).Format(fmt)
}

// dfault checks whether `given` is set, and returns default if not set.
//
// This returns `d` if `given` appears not to be set, and `given` otherwise.
//
// For numeric types 0 is unset.
// For strings, maps, arrays, and slices, len() = 0 is considered unset.
// For bool, false is unset.
// Structs are never considered unset.
//
// For everything else, including pointers, a nil value is unset.
func dfault(d interface{}, given ...interface{}) interface{} {

	if empty(given) || empty(given[0]) {
		return d
	}
	return given[0]
}

// empty returns true if the given value has the zero value for its type.
func empty(given interface{}) bool {
	g := reflect.ValueOf(given)
	if !g.IsValid() {
		return true
	}

	// Basically adapted from text/template.isTrue
	switch g.Kind() {
	default:
		return g.IsNil()
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return g.Len() == 0
	case reflect.Bool:
		return !g.Bool()
	case reflect.Complex64, reflect.Complex128:
		return g.Complex() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return g.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return g.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return g.Float() == 0
	case reflect.Struct:
		return false
	}
}
