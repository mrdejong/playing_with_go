package validate

import (
	"unicode"
)

// Errors is a map holding all the possible errors that may
// occur during validation.
type Errors []string

// Any return true if there is any error.
func (e Errors) Any() bool {
	return len(e) > 0
}

// Add adds an error for a specific field
func (e Errors) Add(msg string) {
	e = append(e, msg)
}

// Get returns all the errors for the given field.
func (e Errors) Get(idx int) string {
	return e[idx]
}

// Schema represents a validation schema.
type Schema struct {
	FieldName string
	Rules     []RuleSet
}

// Rules is a function that takes any amount of RuleSets
func Rules(rules ...RuleSet) []RuleSet {
	ruleSets := make([]RuleSet, len(rules))
	for i := range ruleSets {
		ruleSets[i] = rules[i]
	}
	return ruleSets
}

// Validate validates data based on the given Schema.
func Validate(data any, raw string, fields Schema) (Errors, bool) {
	errors := Errors{}
	return validate(data, raw, fields, errors)
}

func validate(data any, raw string, schema Schema, errors Errors) (Errors, bool) {
	ok := true
	// Uppercase the field name so we never check un-exported fields.
	// But we need to watch out for member fields that are uppercased by
	// the user. For example (URL, ID, ...)
	if !isUppercase(schema.FieldName) {
		schema.FieldName = string(unicode.ToUpper(rune(schema.FieldName[0]))) + schema.FieldName[1:]
	}

	for _, set := range schema.Rules {
		set.FieldValue = data
		set.FieldRawValue = raw
		set.FieldName = schema.FieldName
		if !set.ValidateFunc(set) {
			ok = false
			msg := set.MessageFunc(set)
			if len(set.ErrorMessage) > 0 {
				msg = set.ErrorMessage
			}
			errors = append(errors, msg)
		}
	}
	return errors, ok
}

// func getFieldAndTagByName(v any, name string) any {
// 	val := reflect.ValueOf(v)
// 	if val.Kind() == reflect.Ptr {
// 		val = val.Elem()
// 	}
// 	if val.Kind() != reflect.Struct {
// 		return nil
// 	}
// 	fieldVal := val.FieldByName(name)
// 	if !fieldVal.IsValid() {
// 		return nil
// 	}
// 	return fieldVal.Interface()
// }

func isUppercase(s string) bool {
	for _, ch := range s {
		if !unicode.IsUpper(rune(ch)) {
			return false
		}
	}
	return true
}
