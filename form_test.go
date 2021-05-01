package main

import (
	"reflect"
	"testing"
)

func TestForm_ParseCode(t *testing.T) {
	type fields struct {
		prefix string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Input
		wantErr bool
	}{
		{
			name:    "TestForm_ParseCode: #normal",
			fields:  fields{
				prefix: "_",
			},
			args:    args{
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
				},
				{
					Label: "Is Admin",
					Value: "true",
					Type:  "checkbox",
				},
				{
					Label: "Age",
					Value: "20",
					Type:  "number",
				},
			},
			wantErr: false,
		},

		{
			name:    "TestForm_ParseCode: #should match with the prefix",
			fields:  fields{
				prefix: "_",
			},
			args:    args{
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
				},
			},
			wantErr: false,
		},

		{
			name:    "TestForm_ParseCode: #with default type",
			fields:  fields{
				prefix: "_",
			},
			args:    args{
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
				},
			},
			wantErr: false,
		},

		{
			name:    "TestForm_ParseCode: #error",
			fields:  fields{
				prefix: "_",
			},
			args:    args{
				code: `
$_name = [
	'label' => 'Name',
	'value' => 'billyct',
`,
			},
			want: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Form{
				prefix: tt.fields.prefix,
			}
			got, err := f.ParseCode(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
