// Package srv
package srv

import (
	"awesome-go/pkgs/validate"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type fieldTag struct {
	Name  string
	Rules string
}

func parseFieldTag(sf reflect.StructField) (fieldTag, bool, error) {
	tv := sf.Tag.Get("field")
	if tv == "" {
		return fieldTag{}, false, nil
	}
	parts := strings.SplitN(tv, ":", 2)
	name := strings.TrimSpace(parts[0])
	if name == "" {
		return fieldTag{}, false, fmt.Errorf("invalid field tag on %s (missing name", sf.Name)
	}
	rules := ""
	if len(parts) == 2 {
		rules = strings.TrimSpace(parts[1])
	}
	return fieldTag{Name: name, Rules: rules}, true, nil
}

func buildSchema(t fieldTag) validate.Schema {
	return validate.Schema{
		FieldName: t.Name,
		Rules:     parseRules(t.Rules),
	}
}

var dateLayouts = []string{
	"2006-01-02",
	time.RFC3339,
	"02-01-2006",
	"02 Jan 2006",
	"Jan 2, 2006",
}

func parseDateArg(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	for _, l := range dateLayouts {
		if t, err := time.Parse(l, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func parseRules(rules string) []validate.RuleSet {
	if rules == "" {
		return nil
	}
	var set []validate.RuleSet
	for rule := range strings.SplitSeq(rules, ",") {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		ruleArg := strings.Split(rule, "=")
		key := ruleArg[0]
		arg := ""
		if len(ruleArg) == 2 {
			arg = ruleArg[1]
		}

		switch key {
		case "required":
			set = append(set, validate.Required)
		case "email":
			set = append(set, validate.Email)
		case "min":
			if n, err := strconv.Atoi(arg); err == nil {
				set = append(set, validate.Min(n))
			}
		case "max":
			if n, err := strconv.Atoi(arg); err == nil {
				set = append(set, validate.Max(n))
			}
		case "gte":
			if n, err := strconv.Atoi(arg); err == nil {
				set = append(set, validate.GTE(n))
			}
		case "age":
			if n, err := strconv.Atoi(arg); err == nil {
				set = append(set, validate.Age(n))
			}
		case "before":
			if t, ok := parseDateArg(arg); ok {
				set = append(set, validate.TimeBefore(t))
			}
		case "after":
			if t, ok := parseDateArg(arg); ok {
				set = append(set, validate.TimeAfter(t))
			}
		default:
		}
	}
	return set
}

func mustPtrToStruct(v any) (rv reflect.Value, st reflect.Type, err error) {
	rt := reflect.TypeOf(v)
	if rt == nil || rt.Kind() != reflect.Ptr {
		return reflect.Value{}, nil, fmt.Errorf("not a pointer")
	}
	rv = reflect.ValueOf(v).Elem()
	if rv.Kind() != reflect.Struct {
		return reflect.Value{}, nil, fmt.Errorf("not a pointer to struct")
	}
	return rv, rv.Type(), nil
}

func assignValue(dst reflect.Value, bound any) error {
	bv := reflect.ValueOf(bound)

	if bv.Type().AssignableTo(dst.Type()) {
		dst.Set(bv)
		return nil
	}

	if bv.Kind() == reflect.Ptr && bv.Elem().Type().AssignableTo(dst.Type()) {
		dst.Set(bv.Elem())
		return nil
	}

	if dst.Kind() == reflect.Ptr && bv.Type().AssignableTo(dst.Type().Elem()) {
		p := reflect.New(dst.Type().Elem())
		p.Elem().Set(bv)
		dst.Set(p)
		return nil
	}

	return fmt.Errorf("cannot assign %s to %s", bv.Type(), dst.Type())
}

func walkStruct(rv reflect.Value, st reflect.Type, visit func(dst reflect.Value, sf reflect.StructField)) {
	for i := range rv.NumField() {
		fv := rv.Field(i)
		sf := st.Field(i)

		if fv.Kind() == reflect.Struct && sf.Anonymous {
			walkStruct(fv, fv.Type(), visit)
			continue
		}
		visit(fv, sf)
	}
}

type multiErr []error

func (m multiErr) Error() string {
	var b strings.Builder
	for i, e := range m {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(e.Error())
	}
	return b.String()
}

func joinErrs(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	return multiErr(errs)
}

var (
	binderMu sync.RWMutex
	binders  = map[reflect.Type]func(raw string, schema validate.Schema) any{}
)

func RegisterBinder(t reflect.Type, fn func(raw string, schema validate.Schema) any) {
	binderMu.Lock()
	binders[t] = fn
	binderMu.Unlock()
}

func getBinder(t reflect.Type) (func(string, validate.Schema) any, bool) {
	binderMu.RLock()
	fn, ok := binders[t]
	binderMu.RUnlock()
	return fn, ok
}

func hasBinder(t reflect.Type) bool {
	_, ok := getBinder(t)
	return ok
}

var timeType = reflect.TypeOf(time.Time{})

func walkStructDeep(rv reflect.Value, st reflect.Type, visit func(dst reflect.Value, sf reflect.StructField)) {
	for i := range rv.NumField() {
		fv := rv.Field(i)
		sf := st.Field(i)

		if fv.Kind() == reflect.Struct && !hasBinder(fv.Type()) && fv.Type() != timeType {
			walkStructDeep(fv, fv.Type(), visit)
			continue
		}

		visit(fv, sf)
	}
}

func createField(fieldType reflect.Type, raw string, schema validate.Schema) any {
	isPtr := false
	if fieldType.Kind() == reflect.Ptr {
		isPtr = true
		fieldType = fieldType.Elem()
	}

	fn, ok := getBinder(fieldType)
	if !ok {
		return nil
	}

	out := fn(raw, schema)
	if out == nil {
		return nil
	}

	if isPtr {
		rv := reflect.New(reflect.TypeOf(out))
		rv.Elem().Set(reflect.ValueOf(out))
		return rv.Interface()
	}
	return out
}
