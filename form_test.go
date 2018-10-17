package form

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/lovego/xiaomei/utils"
)

type SubData struct {
	Field4 bool   `form:"field4"`
	Field5 bool   `form:"field5,default=true"`
	Field6 bool   `form:"field6"`
	Field7 bool   `form:"field7"`
	Field8 string `form:"field8,default=sdf"`
	Field9 []string `form:"field9"`
}

type QueryData struct {
	Field1 string `form:"field1"`
	Field2 bool   `form:"field2"`
	Field3 int    `form:"field3,default=5"`

	*SubData
}

func TestFormMapping1(t *testing.T) {
	qd := new(QueryData)
	q := url.Values{}
	err := Bind(q, qd)
	if err != nil {
		t.Fatal(err.Error())
	}
	expected := &QueryData{
		Field3: 5,
		SubData: &SubData{
			Field5: true,
			Field8: "sdf",
		},
	}
	utils.PrintJson(qd)
	if !reflect.DeepEqual(qd, expected) {
		t.Fatalf("unexpected queryData")
	}
}

func TestFormMapping2(t *testing.T) {
	qd := new(QueryData)
	qd.SubData = new(SubData)
	qd.Field9 = []string{"hehe"}
	q := url.Values{
		"field1": {"sdf"},
		"field2": {""},
		"field4": {"true"},
		"field5": {"false"},
		"field6": {"0"},
		"field7": {"1"},
		"field8": {"haha"},
		"field9": {"haha", "hehe"},
	}

	err := Bind(q, qd)
	if err != nil {
		t.Fatal(err.Error())
	}
	expected := &QueryData{
		Field1: "sdf",
		Field3: 5,
		SubData: &SubData{
			Field4: true,
			Field7: true,
			Field8: "haha",
			Field9: []string{"haha", "hehe"},
		},
	}
	if !reflect.DeepEqual(expected, qd) {
		utils.PrintJson(qd)
		t.Fatalf("unexpected queryData")
	}
}


