package main

import "encoding/json"

type Input struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Type  string `json:"type"`
	Name  string `json:"name"`
}

func NewInput() *Input {
	return &Input{
		Type: "text",
	}
}

func NewInputWithMap(m map[string]interface{}) *Input {
	i := NewInput()

	for k, v := range m {
		i.Set(k, v.(string))
	}

	return i
}

func (i *Input) Set(k, v string) {
	switch k {
	case "label":
		i.Label = v
	case "value":
		i.Value = v
	case "type":
		i.Type = v
	case "name":
		i.Name = v
	}
}

func (i *Input) IsEmpty() bool {
	return *NewInput() == *i
}

func (i *Input) ToMap() (map[string]interface{}, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	err = json.Unmarshal(data, &res)
	return res, err
}
