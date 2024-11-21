package dbmap

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"time"
)

var (
	ErrInvalidStruct = errors.New("result must be a struct to be decoded into")
)

func Decode(val map[string]any, result any) error {
	if result == nil {
		slog.Error("dbmap", "message", "result is nil")
		return ErrInvalidStruct
	}
	typ := reflect.TypeOf(result)
	slog.Debug("dbmap", "type", typ.Kind())
	slog.Debug("dbmap", "val", val)
	if typ.Kind() != reflect.Ptr {
		return ErrInvalidStruct
	}
	typ = typ.Elem()
	if typ.Kind() != reflect.Struct {
		return ErrInvalidStruct
	}
	valOfResult := reflect.ValueOf(result).Elem()
	if !valOfResult.IsValid() {
		slog.Error("dbmap", "message", "valOfResult is invalid")
		return ErrInvalidStruct
	}

	for i := 0; i < typ.NumField(); i++ {
		fld := typ.Field(i)
		key, ok := fld.Tag.Lookup("dbmap")
		if !ok {
			slog.Debug("dbmap", "message", fmt.Sprintf("no mapping key for %s", fld.Name))
			continue
		}
		mapVal, exists := val[key]
		valField := valOfResult.Field(i)
		if !valField.CanSet() {
			continue
		}
		switch fld.Type {
		case reflect.TypeOf(&NilString{}):
			slog.Debug("dbmap", "key", key, "type", "NilString")
			strVal, isValidString := mapVal.(string)

			var nilVal *NilString
			if !isValidString && exists {
				slog.Error("dbmap", "message", "value is not a valid string", "key", key, "val", mapVal)
				return ErrInvalidStruct
			}
			nilVal = &NilString{
				Val:   strVal, // is a string with value of empty string if value is not found
				Empty: !exists,
			}
			slog.Debug("dbmap", "key", key, "type", "NilString", "OGVal", mapVal, "resultVal", nilVal, "exists", exists)
			valField.Set(reflect.ValueOf(nilVal))

		case reflect.TypeOf(&NilInt{}):
			slog.Debug("dbmap", "key", key, "type", "NilInt")
			intVal, isValidInt := mapVal.(int)

			var nilVal *NilInt
			if !isValidInt && exists {
				slog.Error("dbmap", "message", "value is not a valid int", "key", key, "val", mapVal)
				return ErrInvalidStruct
			}
			nilVal = &NilInt{
				Val:   intVal, // is an int with value of 0 if value is not found
				Empty: !exists,
			}
			slog.Debug("dbmap", "key", key, "type", "NilInt", "OGVal", mapVal, "resultVal", nilVal, "exists", exists)
			valField.Set(reflect.ValueOf(nilVal))
		case reflect.TypeOf(&NilTime{}):
			slog.Debug("dbmap", "key", key, "type", "NilTime")
			timeVal, isValidTime := mapVal.(time.Time)

			var nilVal *NilTime
			if !isValidTime && exists {
				slog.Error("dbmap", "message", "value is not a valid time", "key", key, "val", mapVal)
				return ErrInvalidStruct
			}
			nilVal = &NilTime{
				Val:   timeVal, // is a time with value of zero time if value is not found
				Empty: !exists,
			}
			slog.Debug("dbmap", "key", key, "type", "NilTime", "OGVal", mapVal, "resultVal", nilVal, "exists", exists)
			valField.Set(reflect.ValueOf(nilVal))
		default:
			if exists {
				valField.Set(reflect.ValueOf(mapVal))
			}
			slog.Debug("dbmap", "key", key, "type", "not any of the niltype")
		}
	}
	return nil
}
