package utils

import (
	"gopkg.in/go-playground/validator.v8"
	"reflect"
)

func NameValidator(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	// Require the type is String and contains 1 to 10 characters.
	if fieldKind != reflect.String {
		return false
	}
	value := field.String()
	if value == "" || len(value) > 10 {
		return false
	}
	return true
}

func MobileValidator(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	// Require the type is String and contains 0 or 11 numeric characters.
	if fieldKind != reflect.String {
		return false
	}
	value := field.String()

	if len(value) == 11 {
		return validator.IsNumeric(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param)
	} else if len(value) == 0 {
		return true
	}
	return false
}

func EmailValidator(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	// Require the type is String and length less than 100, and Email format.
	if fieldKind != reflect.String {
		return false
	}
	value := field.String()
	if len(value) > 100 {
		return false
	}
	return validator.IsEmail(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param)
}

func PasswordValidator(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	// Require the type is String, and length between 8 to 12.
	if fieldKind != reflect.String {
		return false
	}
	value := field.String()
	if len(value) < 8 || len(value) > 12 {
		return false
	}
	return true
}

func GenderValidator(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	// Require the type is Int, and value is -1 or 0 or 1
	if fieldKind != reflect.Int {
		return false
	}
	value := field.Int()
	if value < -1 || value > 1 {
		return false
	}
	return true
}
