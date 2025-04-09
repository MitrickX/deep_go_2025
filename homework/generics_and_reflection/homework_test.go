package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

const (
	LexemProperties = "properties"
	LexemOmitEmpty  = "omitempty"
)

func parseTag(tag reflect.StructTag) (fieldName string, omitempty bool, ok bool) {
	meta, ok := tag.Lookup(LexemProperties)
	if !ok {
		return
	}

	parts := strings.Split(meta, ",")

	fieldName = strings.TrimSpace(parts[0])
	if len(fieldName) == 0 {
		ok = false
	}

	if len(parts) == 1 {
		return
	}

	omitempty = strings.TrimSpace(parts[1]) == LexemOmitEmpty
	return
}

func Serialize(person Person) string {
	personType := reflect.TypeOf(person)
	personValue := reflect.ValueOf(person)

	bs := strings.Builder{}

	n := personType.NumField()
	for i := 0; i < n; i++ {
		personField := personType.Field(i)
		fieldName, omitempty, ok := parseTag(personField.Tag)
		if !ok {
			continue
		}

		fieldValue := personValue.Field(i)
		if omitempty && fieldValue.IsZero() {
			continue
		}

		if bs.Len() > 0 {
			bs.WriteByte('\n')
		}

		bs.WriteString(fieldName)
		bs.WriteByte('=')
		bs.WriteString(fmt.Sprintf("%v", fieldValue))
	}

	return bs.String()
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
