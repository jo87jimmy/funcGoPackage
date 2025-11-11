package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
)

// 從 echo context 取得 query 或 form 的參數
func RequestParams(c echo.Context, params ...string) map[string]string {
	pargs := map[string]string{}

	for _, arg := range params {
		key := fmt.Sprintf("params[%s]", arg)
		if val := c.QueryParam(key); val != "" {
			pargs[arg] = val
		} else if val := c.FormValue(key); val != "" {
			pargs[arg] = val
		}
	}

	return pargs
}

// 將 struct 或 map 轉為 ORM Where 條件
func OrmWhereParams(input any, SQLTable any) map[string]any {
	output := make(map[string]any, 100)
	v := reflect.ValueOf(input)
	t := v.Type()
	s := reflect.ValueOf(SQLTable)
	sqlType := s.Type()
	switch v.Kind() {
	case reflect.Map:
		for _, key := range v.MapKeys() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			val := v.MapIndex(key).Interface()
			match := func(fieldName string) bool {
				return strings.EqualFold(fieldName, keyStr)
			}

			if field, ok := sqlType.FieldByNameFunc(match); ok {
				column, _ := field.Tag.Lookup("db")
				if column != "" {
					output[column] = val
				}
			}
		}

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldValue := v.Field(i)
			if !fieldValue.CanInterface() {
				continue
			}
			if fieldValue.Kind() == reflect.String && fieldValue.Len() == 0 {
				continue
			}

			fieldName := t.Field(i).Name
			match := func(fieldName2 string) bool {
				return strings.EqualFold(fieldName2, fieldName)
			}

			if field, ok := sqlType.FieldByNameFunc(match); ok {
				column, _ := field.Tag.Lookup("db")
				if column != "" {
					output[column] = fieldValue.Interface()
				}
			}
		}
	}

	return output
}
