package main

import (
	"reflect"
	"testing"
)

func TestInput_IsEmpty(t *testing.T) {
	type fields struct {
		Label string
		Value string
		Type  string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "TestInput_IsEmpty: all empty",
			fields: fields{
				Label: "",
				Value: "",
				Type:  "",
			},
			want: false,
		},
		{
			name: "TestInput_IsEmpty: with default type",
			fields: fields{
				Label: "",
				Value: "",
				Type:  "text",
			},
			want: true,
		},
		{
			name: "TestInput_IsEmpty: with label",
			fields: fields{
				Label: "label",
				Value: "",
				Type:  "text",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Input{
				Label: tt.fields.Label,
				Value: tt.fields.Value,
				Type:  tt.fields.Type,
			}
			if got := i.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInput_Set(t *testing.T) {
	type args struct {
		k string
		v string
	}
	tests := []struct {
		name string
		args args
		want *Input
	}{
		{
			name: "TestInput_Set: should be correctly",
			args: args{
				k: "label",
				v: "label",
			},
			want: &Input{
				Label: "label",
				Value: "",
				Type:  "",
			},
		},

		{
			name: "TestInput_Set: should be correctly",
			args: args{
				k: "other",
				v: "label",
			},
			want: &Input{
				Label: "",
				Value: "",
				Type:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Input{}
			if i.Set(tt.args.k, tt.args.v); *i != *tt.want {
				t.Errorf("Set() = %v, want %v", i, tt.want)
			}
		})
	}
}

func TestNewInput(t *testing.T) {
	tests := []struct {
		name string
		want *Input
	}{
		{
			name: "TestNewInput",
			want: &Input{
				Type: "text",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInput(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
