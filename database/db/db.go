package db

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/iancoleman/strcase"
	"reflect"
	"time"
)

func skipSpecialField(field *structs.Field) bool {
	fieldName := strcase.ToSnake(field.Name())
	asName := fieldName == "created_at" || fieldName == "deleted_at" || fieldName == "id" || fieldName == "updated_at"
	asTag := field.Tag("q") == "_"

	return asName || asTag
}

func D(T any) (q string) {
	t := structs.New(T)
	q = fmt.Sprintf("delete from %s  ", strcase.ToSnake(t.Name()))
	return q + fmt.Sprintf(" where id = %d", t.Field("Id").Value())
}

func U(T any) (q string) {
	t := structs.New(T)
	q = fmt.Sprintf("update %s set ", strcase.ToSnake(t.Name()))
	v := t.Fields()
	for i := 0; i < len(v); i++ {
		if skipSpecialField(v[i]) {
			continue
		}

		switch v[i].Kind() {
		case reflect.String:
			q += fmt.Sprintf("%s = '%s',", strcase.ToSnake(v[i].Name()), v[i].Value())
			break
		case reflect.Struct:
			if w, ok := v[i].Value().(time.Time); ok {
				q += fmt.Sprintf("%s = '%s',", strcase.ToSnake(v[i].Name()), w.Format(time.DateTime))
			}
			break
		default:
			q += fmt.Sprintf("%s = %v,", strcase.ToSnake(v[i].Name()), v[i].Value())
			break
		}
	}

	return q[:len(q)-1] + fmt.Sprintf(" where id = %d", t.Field("Id").Value())
}

func I(T any) (q string) {
	t := structs.New(T)
	q = fmt.Sprintf("insert into %s (", strcase.ToSnake(t.Name()))
	v := t.Fields()
	for i := 0; i < len(v); i++ {
		if skipSpecialField(v[i]) {
			continue
		}
		q += fmt.Sprintf("%s,", strcase.ToSnake(v[i].Name()))
	}
	q = q[:len(q)-1] + ") VALUES ("
	for i := 0; i < len(v); i++ {
		if skipSpecialField(v[i]) {
			continue
		}
		switch v[i].Kind() {
		case reflect.String:
			q += fmt.Sprintf("'%s',", v[i].Value())
			break
		case reflect.Struct:
			if w, ok := v[i].Value().(time.Time); ok {
				q += fmt.Sprintf("'%s',", w.Format(time.DateTime))
			}
			break
		default:
			q += fmt.Sprintf("%v,", v[i].Value())
			break
		}
	}

	q = q[:len(q)-1] + ")"

	return q
}
