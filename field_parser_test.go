package form

import (
    "strings"
    "testing"
    "net/url"
    "reflect"
    "fmt"
)

type Field1 map[string]string
type Field2 []string
type Field3 [2]string
type Field4 []string

type CustomFieldQuery struct {
    Field1  `form:"field1"`
    Field2 `form:"field2"`
    Field3 *Field3 `form:"field3"`
}

type WrongQueryData struct {
    Field4 `form:"field4"`
}

func(f1 Field1) FieldParse(input string)  (interface{}, error){
    out := make(Field1)
    tups := strings.Split(input, ",")
    for _, v := range tups{
        stps := strings.Split(v, ":")
        out[strings.TrimSpace(stps[0])] = strings.TrimSpace(stps[1])
    }
    return out, nil
}

func (f2 Field2) FieldParse(input string)  (interface{}, error){
    return strings.Split(input, ","), nil
}

func (f3 *Field3) FieldParse(input string)  (interface{}, error){
    out := Field3{}
    sp := strings.Split(input, ",")
    out[0], out[1] = sp[0], sp[1]
    return out, nil
}

func (f4 Field4) FieldParse(input string)(interface{}, error){
    return [2]string{"sdf", "aaa"}, nil
}


func TestCustomField(t *testing.T){
    query := url.Values{
        "field1": {"A:a, B:b"},
        "field2": {"2,3,4"},
        "field3": {"2,3"},
    }
    queryData := new(CustomFieldQuery)
    err := Bind(query, queryData)
    if err != nil{
        t.Fatal(err)
    }
    expField1 := Field1{"A": "a", "B": "b"}
    expField2 := Field2{"2", "3", "4"}
    expField3 := Field3{"2", "3"}
    if !(reflect.DeepEqual(expField1, queryData.Field1) &&
        reflect.DeepEqual(expField2, queryData.Field2) &&
        reflect.DeepEqual(expField3, *queryData.Field3)){
        fmt.Printf("%#v\n", queryData)
        t.Fatal("test failed")
    }
}



func TestWrongParser(t *testing.T){

    query := url.Values{
        "field4": {"adf, haha"},
    }
    queryData := new(WrongQueryData)
    err := Bind(query, queryData)
    if err == nil || err.Error() != "field: Field4 wrong parsed type: [2]string"{
        t.Fatal(err)
    }else{
        t.Log(err)
    }

}


