// Package srv server package
package srv

import (
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type Request struct {
	c *fiber.Ctx
}

type ErrorGetter interface{ AnyErrors() bool }

func HasValidationErrors(v any) bool {
	rv, st, err := mustPtrToStruct(v)
	if err != nil {
		return true
	} // treat shape issues as “has problems”
	found := false

	walkStructDeep(rv, st, func(dst reflect.Value, _ reflect.StructField) {
		if found {
			return
		}
		val := dst
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				return
			}
			val = val.Elem()
		}
		if !val.IsValid() {
			return
		}

		ig := val.Interface()
		if eg, ok := ig.(ErrorGetter); ok && eg.AnyErrors() {
			found = true
			return
		}
	})
	return found
}

func extractErrorGetter(bound any) (ErrorGetter, bool) {
	// handle both value and pointer results from binders
	v := reflect.ValueOf(bound)
	if eg, ok := bound.(ErrorGetter); ok {
		return eg, true
	}
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		if eg, ok := v.Elem().Interface().(ErrorGetter); ok {
			return eg, true
		}
	}
	return nil, false
}

func (*Server) ParseFields(v any, getter func(name string, defaultVal ...string) string) (hasErrors bool) {
	rv, st, err := mustPtrToStruct(v)
	if err != nil {
		panic(err)
	}

	var errs []error
	walkStructDeep(rv, st, func(dst reflect.Value, sf reflect.StructField) {
		tag, ok, terr := parseFieldTag(sf)
		if terr != nil {
			errs = append(errs, terr)
			return
		}
		if !ok {
			return
		}

		schema := buildSchema(tag)
		raw := getter(tag.Name, "")
		bound := createField(sf.Type, raw, schema)
		if bound == nil {
			return
		}

		if eg, ok := extractErrorGetter(bound); ok && eg.AnyErrors() {
			hasErrors = true
		}
		if aerr := assignValue(dst, bound); aerr != nil {
			errs = append(errs, fmt.Errorf("%s: %w", sf.Name, aerr))
		}
	})
	if es := joinErrs(errs); es != nil {
		panic(es)
	}
	return hasErrors
}
