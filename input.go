package main

type Input struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

func NewInput() *Input {
	return &Input{
		Type: "text",
	}
}

func (i *Input) Set(k, v []byte)  {
	sk, sv := string(k), string(v)

	switch sk {
	case "label":
		i.Label = sv
	case "value":
		i.Value = sv
	case "type":
		i.Type = sv
	}
}

func (i *Input) IsEmpty() bool {
	return *NewInput() == *i
}
