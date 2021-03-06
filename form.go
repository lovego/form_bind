// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package form

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
    "fmt"
)

func Bind(q url.Values, obj interface{}) error {
	if err := mapForm(obj, q); err != nil {
		return err
	}
	return nil
}

type InvalidUnmarshalError struct {
    Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
    if e.Type == nil {
        return "form: Unmarshal(nil)"
    }

    if e.Type.Kind() != reflect.Ptr {
        return "form: Unmarshal(non-pointer " + e.Type.String() + ")"
    }
    return "form: Unmarshal(nil " + e.Type.String() + ")"
}

type FieldParser interface {
    FieldParse(string) (interface{}, error)
}

// 逗号之间最好不要有空格
func mapForm(ptr interface{}, form map[string][]string) error {
	typ := reflect.TypeOf(ptr).Elem()
	val := reflect.ValueOf(ptr).Elem()
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		if !structField.CanSet() {
			continue
		}

		structFieldKind := structField.Kind()
		inputFieldName := typeField.Tag.Get("form")
		inputFieldNameList := strings.Split(inputFieldName, ",")
		inputFieldName = inputFieldNameList[0]
		var defaultValue string
		if len(inputFieldNameList) > 1 {
			defaultList := strings.SplitN(strings.TrimSpace(inputFieldNameList[1]), "=", 2)
			if defaultList[0] == "default" {
				defaultValue = defaultList[1]
			}
		}
		if inputFieldName == "" {
			inputFieldName = typeField.Name

			// if "form" tag is nil, we inspect if the field is a struct or struct pointer.
			// this would not make sense for JSON parsing but it does for a form
			// since data is flatten
			if structFieldKind == reflect.Ptr {
				structField = getStructFieldByPtr(structField)
				structFieldKind = structField.Kind()
			}
			if structFieldKind != reflect.Struct{
			    return errors.New("if no form tag provided, it must be a struct or its pointer")
            }
			if structFieldKind == reflect.Struct {
				err := mapForm(structField.Addr().Interface(), form)
				if err != nil {
					return err
				}
				continue
			}
		}
		inputValue, exists := form[inputFieldName]

		if !exists {
			if defaultValue == "" {
				continue
			}
			inputValue = make([]string, 1)
			inputValue[0] = defaultValue
		}
		if fieldInterf, isFieldBinder := structField.Interface().(FieldParser); isFieldBinder{
            if structFieldKind == reflect.Ptr {
                structField = getStructFieldByPtr(structField)
                structFieldKind = structField.Kind()
            }
            parsed, err := fieldInterf.FieldParse(inputValue[0])
            if err != nil{
                return err
            }
            if reflect.TypeOf(parsed).Kind() != structFieldKind{
                return fmt.Errorf("field: %s wrong parsed type: %T", typeField.Name, parsed)
            }
            structField.Set(reflect.ValueOf(parsed))
            continue
        }
        if structFieldKind == reflect.Slice && len(inputValue) > 0 {
            setSliceField(structField, inputValue)
            continue
        }
		if _, isTime := structField.Interface().(time.Time); isTime {
			if err := setTimeField(inputValue[0], typeField, structField); err != nil {
				return err
			}
			continue
		}
		if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil {
			return err
		}
	}
	return nil
}

func getStructFieldByPtr(structField reflect.Value) reflect.Value{

    if !structField.Elem().IsValid() {
        structField.Set(reflect.New(structField.Type().Elem()))
    }
    structField = structField.Elem()
    return structField
}

func setSliceField(structField reflect.Value, inputValue []string) error{
    elemsNum := len(inputValue)
    sliceOf := structField.Type().Elem().Kind()
    slice := reflect.MakeSlice(structField.Type(), elemsNum, elemsNum)
    for i := 0; i < elemsNum; i++ {
        if err := setWithProperType(sliceOf, inputValue[i], slice.Index(i)); err != nil {
            return err
        }
    }
    structField.Set(slice)
    return nil
}

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	switch valueKind {
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32:
		return setFloatField(val, 32, structField)
	case reflect.Float64:
		return setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	case reflect.Ptr:
		if !structField.Elem().IsValid() {
			structField.Set(reflect.New(structField.Type().Elem()))
		}
		structFieldElem := structField.Elem()
		return setWithProperType(structFieldElem.Kind(), val, structFieldElem)
	default:
		return errors.New("Unknown type")
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}

	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}
