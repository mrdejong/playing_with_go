package srv

import (
	"awesome-go/pkgs/validate"
	"time"
)

type FieldBase struct {
	Name     string
	RawValue string
	Errors   validate.Errors
}

func (fb FieldBase) AnyErrors() bool { return fb.Errors.Any() }

type StringField struct {
	FieldBase
	Value string
}

type NumberField struct {
	FieldBase
	Value int
}

type DateField struct {
	FieldBase
	Value time.Time
}
