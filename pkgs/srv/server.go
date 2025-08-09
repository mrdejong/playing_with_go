package srv

import (
	"awesome-go/pkgs/validate"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	numberFieldType = reflect.TypeOf(NumberField{})
	stringFieldType = reflect.TypeOf(StringField{})
	dateFieldType   = reflect.TypeOf(DateField{})
)

type Server struct {
	// Define server fields here
	Router *fiber.App
}

func New() *Server {
	initializeBinders()
	return &Server{
		Router: fiber.New(),
	}
}

func initializeBinders() {
	RegisterBinder(numberFieldType, func(raw string, schema validate.Schema) any {
		iv, err := strconv.Atoi(raw)
		if err != nil {
			return NumberField{FieldBase: FieldBase{Name: schema.FieldName, RawValue: raw, Errors: validate.Errors{"must be a numeric value"}}, Value: 0}
		}
		errs, _ := validate.Validate(iv, raw, schema)
		return NumberField{FieldBase: FieldBase{Name: schema.FieldName, RawValue: raw, Errors: errs}, Value: iv}
	})

	RegisterBinder(stringFieldType, func(raw string, schema validate.Schema) any {
		errs, _ := validate.Validate(raw, raw, schema)
		return StringField{FieldBase: FieldBase{Name: schema.FieldName, RawValue: raw, Errors: errs}, Value: raw}
	})

	RegisterBinder(dateFieldType, func(raw string, schema validate.Schema) any {
		df := DateField{FieldBase: FieldBase{Name: schema.FieldName, RawValue: raw}}

		layouts := []string{
			"2006-01-02",
			"02-01-2006",
			time.RFC3339,
			"02 Jan 2006",
			"Jan 2, 2006",
		}

		raw = strings.TrimSpace(raw)
		if raw == "" {
			errs, _ := validate.Validate(time.Time{}, raw, schema)
			df.Errors = errs
			return df
		}

		var t time.Time
		var perr error
		for _, l := range layouts {
			if tt, err := time.Parse(l, raw); err == nil {
				t = tt
				perr = nil
				break
			} else {
				perr = err
			}
		}
		if !t.IsZero() {
			df.Value = t
			errs, _ := validate.Validate(t, raw, schema)
			df.Errors = errs
			return df
		}

		df.Errors = append(df.Errors, perr.Error())
		return df
	})
}
