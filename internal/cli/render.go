package cli

import (
	"github.com/rodaine/table"
	"reflect"
)

func RenderEntityTable(entity any) {
	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	tbl := table.New("Field", "Value")
	for i := 0; i < v.NumField(); i++ {
		name := t.Field(i).Name
		value := v.Field(i).Interface()
		tbl.AddRow(name, value)
	}
	tbl.Print()
}

func RenderCollectionTable[T any](collection []T) {
	if len(collection) == 0 {
		return
	}
	for _, item := range collection {
		RenderEntityTable(item)
	}
}
