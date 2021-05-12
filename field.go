package main

type Field struct {
	Name  string
	Value string
	Data  map[string]interface{}
}

func NewField() *Field {
	return &Field{}
}

func NewFieldWithMap(m map[string]interface{}) *Field {
	f := NewField()

	for k, v := range m {
		f.Set(k, v)
	}

	return f
}

func (f *Field) SetName(name string) {
	f.Name = name
}

func (f *Field) GetName() string {
	return f.Name
}

func (f *Field) SetValue(v string) {
	f.Value = v
}

func (f *Field) Set(k string, v interface{}) {
	if k == "" {
		return
	}

	switch k {
	case "value":
		f.Value = v.(string)
	case "name":
		f.Name = v.(string)
	default:
		if f.Data == nil {
			f.Data = make(map[string]interface{})
		}

		f.Data[k] = v
	}
}

func (f *Field) IsEmpty() bool {
	return f.Name == "" && f.Value == "" && f.Data == nil
}

func (f *Field) ToMap() map[string]interface{} {
	res := f.Data
	if res == nil {
		res = make(map[string]interface{})
	}

	res["name"] = f.Name
	res["value"] = f.Value
	return res
}
