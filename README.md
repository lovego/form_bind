# form_bind
动态绑定form到结构体

支持的tag有：`form`, `time_format`, `time_utc`, `time_location`

如果没有`form` tag，那么这个字段必须为结构体或结构体指针

例如
```
struct Query{
    Field1 time.Time `form:"field1" time_format:"2006-01-02T15:04:05Z07:00"` // 默认时间格式：time.RFC3339
    Field2 time.Time `form:"field2" time_utc:"true"` // utc时间
    Field3 time.Time `form:"filed3" time_location:"Local"` // local时间
}
```


支持自定义字段解析，只需要实现 `FieldParse(string) (interface{}, error)`方法即可,

需要注意的是FieldParse返回的值要是原字段类型，否则会返回错误。

详细使用见`field_parser_test.go`