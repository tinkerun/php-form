package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestForm_ParseCode(t *testing.T) {
	type fields struct {
		prefix string
		code   string
	}

	tests := []struct {
		name    string
		fields  fields
		want    []Input
		wantErr bool
	}{
		{
			name: "TestForm_ParseCode: #normal",
			fields: fields{
				prefix: "_",
				code: `
$_name = [
	'label' => 'Name',
	'value' => 'billyct',
	'type'  => 'text',
];
$_isAdmin = [
	'label' => 'Is Admin',
	'value' => true,
	'type'  => 'checkbox',
];
$_age = [
	'label' => 'Age',
	'value' => 20,
	'type'  => 'number',
];
`,
			},
			want: []Input{
				{
					Label: "Name",
					Value: "billyct",
					Type:  "text",
					Name:  "$_name",
				},
				{
					Label: "Is Admin",
					Value: "true",
					Type:  "checkbox",
					Name:  "$_isAdmin",
				},
				{
					Label: "Age",
					Value: "20",
					Type:  "number",
					Name:  "$_age",
				},
			},
			wantErr: false,
		},

		{
			name: "TestForm_ParseCode: #should match with the prefix",
			fields: fields{
				prefix: "_",
				code: `
$_name = [
	'label' => 'Name',
	'value' => 'billyct',
	'type'  => 'text',
];
$isAdmin = [
	'label' => 'Is Admin',
	'value' => true,
	'type'  => 'checkbox',
];
$age = [
	'label' => 'Age',
	'value' => 20,
	'type'  => 'number',
];
`,
			},
			want: []Input{
				{
					Label: "Name",
					Value: "billyct",
					Type:  "text",
					Name:  "$_name",
				},
			},
			wantErr: false,
		},

		{
			name: "TestForm_ParseCode: #with default type",
			fields: fields{
				prefix: "_",
				code: `
$_name = [
	'label' => 'Name',
	'value' => 'billyct',
];
`,
			},
			want: []Input{
				{
					Label: "Name",
					Value: "billyct",
					Type:  "text",
					Name:  "$_name",
				},
			},
			wantErr: false,
		},

		{
			name: "TestForm_ParseCode: #error",
			fields: fields{
				prefix: "_",
				code: `
$_name = [
	'label' => 'Name',
	'value' => 'billyct',
`,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Form{
				prefix: tt.fields.prefix,
				code:   tt.fields.code,
			}
			got, err := f.GenerateInputs()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateInputs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateInputs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForm_GenerateCodeWithInputs(t *testing.T) {
	type fields struct {
		prefix string
		code   string
	}
	type args struct {
		inputs []Input
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "TestForm_GenerateCodeWithInputs",
			fields: fields{
				prefix: "_",
				code: `
$_name = [
	'label' => 'Name',
	'value' => 'billyct',
];
$_isAdmin = [
	'label' => 'Is Admin',
	'value' => true,
];
$_age = [
	'label' => 'Age',
	'value' => 20,
];
`,
			},
			args: args{
				inputs: []Input{
					{
						Value: "hello",
						Name:  "$_name",
					},
					{
						Value: "30",
						Name:  "$_age",
					},
					{
						Value: "false",
						Name:  "$_isAdmin",
					},
				},
			},
			want: `$_name = [
	'label' => 'Name',
	'value' => 'hello',
];
$_isAdmin = [
	'label' => 'Is Admin',
	'value' => false,
];
$_age = [
	'label' => 'Age',
	'value' => 30,
];`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Form{
				prefix: tt.fields.prefix,
				code:   tt.fields.code,
			}
			got, err := f.GenerateCodeWithInputs(tt.args.inputs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateCodeWithInputs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(got, tt.want) {
				t.Errorf("GenerateCodeWithInputs() got = %v, want %v", got, tt.want)
			}
		})
	}
}
