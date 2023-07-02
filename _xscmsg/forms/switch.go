package xscmsg

import (
	"fmt"
	"reflect"
)

type validator interface{ Validate() []string }

var validatorType = reflect.TypeOf((*validator)(nil)).Elem()

var validators = map[reflect.Type]reflect.Type{}

// Register registers a handler type that provides functionality for a
// particular message contents type.  The handler type must be a structure whose
// first member is the message contents type that it provides functionality for.
func Register(handler any) {
	htype := reflect.TypeOf(handler)
	if htype.Kind() != reflect.Struct || htype.NumField() == 0 {
		panic("parameter to Register must be struct with at least one field")
	}
	if htype.Implements(validatorType) {
		validators[htype.Field(0).Type] = htype
		println("registered!")
	}
}

// Validate validates the message contents.  It returns a list of strings
// describing problems with the contents; the list is empty if there are none.
func Validate(msg any) (problems []string) {
	vtype := validators[reflect.TypeOf(msg)]
	if vtype == nil {
		return []string{fmt.Sprintf("The message contents could not be validated because there is no registered validator for %T.", msg)}
	}
	vptr := reflect.New(vtype)
	vptr.Elem().Field(0).Set(reflect.ValueOf(msg))
	return vptr.Interface().(validator).Validate()
}
