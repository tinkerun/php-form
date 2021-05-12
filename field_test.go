package main

import (
	"reflect"
	"testing"
)

func TestNewFieldWithMap(t *testing.T) {
	type args struct {
		m map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want *Field
	}{
		{
			name: "TestNewFieldWithMap: #1",
			args: args{
				m: map[string]interface{}{
					"name": "Name",
				},
			},
			want: &Field{
				Name: "Name",
			},
		},
		{
			name: "TestNewFieldWithMap: #2",
			args: args{
				m: map[string]interface{}{
					"name":   "Name",
					"value":  "default value",
					"other":  "other",
					"other2": "other2",
				},
			},
			want: &Field{
				Name:  "Name",
				Value: "default value",
				Data: map[string]interface{}{
					"other":  "other",
					"other2": "other2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFieldWithMap(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFieldWithMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_Set(t *testing.T) {
	type args struct {
		k string
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want *Field
	}{
		{
			name: "TestField_Set: #1",
			args: args{
				k: "name",
				v: "hello",
			},
			want: &Field{
				Name: "hello",
			},
		},

		{
			name: "TestField_Set: #2",
			args: args{
				k: "value",
				v: "hello",
			},
			want: &Field{
				Value: "hello",
			},
		},

		{
			name: "TestField_Set: #3",
			args: args{
				k: "other",
				v: "other",
			},
			want: &Field{
				Data: map[string]interface{}{
					"other": "other",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewField()

			if f.Set(tt.args.k, tt.args.v); !reflect.DeepEqual(f, tt.want) {
				t.Errorf("Set() = %v, want %v", f, tt.want)
			}
		})
	}
}

func TestField_ToMap(t *testing.T) {
	type fields struct {
		Name  string
		Value string
		Data  map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{
			name: "TestField_ToMap: #1",
			fields: fields{
				Name:  "name1",
				Value: "value1",
				Data:  nil,
			},
			want: map[string]interface{}{
				"name":  "name1",
				"value": "value1",
			},
		},
		{
			name: "TestField_ToMap: #2",
			fields: fields{
				Name:  "name2",
				Value: "value2",
				Data: map[string]interface{}{
					"type":  "type2",
					"label": "label2",
				},
			},
			want: map[string]interface{}{
				"name":  "name2",
				"value": "value2",
				"type":  "type2",
				"label": "label2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Field{
				Name:  tt.fields.Name,
				Value: tt.fields.Value,
				Data:  tt.fields.Data,
			}
			if got := f.ToMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
