package validator

import (
	"log"
	"testing"
)

func TestValidator(t *testing.T) {
	v := New()
	data := map[string]string{
		"int64param": "32",
		"int32param": "3213",
	}
	rules := map[string]string{
		"int64param": "int64",
		"int32param": "int32",
	}
	v.Validate(&data, rules)
	if v.HasErr == true {
		log.Fatal("test fail")
	}

	v = New()
	data = map[string]string{
		"int64param": "32",
		"int32param": "aasd",
	}
	rules = map[string]string{
		"int64param": "int64",
		"int32param": "int32",
	}
	v.Validate(&data, rules)
	if v.HasErr == false {
		log.Fatal("test fail")
	}

	// test datetime
	v = New()
	data = map[string]string{
		"datetime": "2016-09-13 08:33:12s",
	}
	rules = map[string]string{
		"datetime": "datetime",
	}
	v.Validate(&data, rules)
	if v.HasErr == false {
		log.Fatal("datetime test fail")
	}
	log.Println("datetime test success")

	v = New()
	data = map[string]string{
		"datetime": "2016-09-13 08:33:12",
	}
	rules = map[string]string{
		"datetime": "datetime",
	}
	v.Validate(&data, rules)
	if v.HasErr == true {
		log.Fatal("datetime test fail")
	}
	log.Println("datetime test success")

	v = New()
	data = map[string]string{
		"data":  "1",
		"data2": "4",
	}
	rules = map[string]string{
		"data":  "in:1,2,3",
		"data2": "notIn:1,2,3",
	}
	v.Validate(&data, rules)
	if v.HasErr == false {
		log.Println("in not in test success")
	} else {
		log.Fatal("in not in test fail")
	}
}
