package main

import (
	"reflect"
	"testing"
)

func TestForm_Parse(t *testing.T) {
	type fields struct {
		prefix string
		code   string
	}

	tests := []struct {
		name    string
		fields  fields
		want    []Field
		wantErr bool
	}{
		{
			name: "TestForm_Parse: normal",
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
			want: []Field{
				{
					Name:  "$_name",
					Value: "billyct",
					Data: map[string]interface{}{
						"label": "Name",
						"type":  "text",
					},
				},
				{
					Name:  "$_isAdmin",
					Value: "true",
					Data: map[string]interface{}{
						"label": "Is Admin",
						"type":  "checkbox",
					},
				},
				{
					Value: "20",
					Name:  "$_age",
					Data: map[string]interface{}{
						"label": "Age",
						"type":  "number",
					},
				},
			},
			wantErr: false,
		},

		{
			name: "TestForm_Parse: short",
			fields: fields{
				prefix: "_",
				code: `
$_name = 'billyct';
$_isAdmin = true;
$_age = 20;
`,
			},
			want: []Field{
				{
					Name:  "$_name",
					Value: "billyct",
				},
				{
					Name:  "$_isAdmin",
					Value: "true",
				},
				{
					Value: "20",
					Name:  "$_age",
				},
			},
			wantErr: false,
		},

		{
			name: "TestForm_Parse: should match with the prefix",
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
			want: []Field{
				{
					Value: "billyct",
					Name:  "$_name",
					Data: map[string]interface{}{
						"label": "Name",
						"type":  "text",
					},
				},
			},
			wantErr: false,
		},

		{
			name: "TestForm_Parse: with default type",
			fields: fields{
				prefix: "_",
				code: `
$_name = [
	'label' => 'Name',
	'value' => 'billyct',
];
`,
			},
			want: []Field{
				{
					Value: "billyct",
					Name:  "$_name",
					Data: map[string]interface{}{
						"label": "Name",
					},
				},
			},
			wantErr: false,
		},

		{
			name: "TestForm_Parse: error",
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
		{
			name: "TestForm_Parse: case 1",
			fields: fields{
				prefix: "form",
				code: `
<?php

use App\Models\User;
use Illuminate\Support\Facades\DB;

$form_email = [
	'label' => 'Email',
	'value' => '',
	'type' => 'text',
];

// Customer Support
// when a user does not receive a password reset email

$user = User::where('email', $form_email['value'])->first();

$user->password = bcrypt('your-new-secure-password');

$user->save();

$user;
`,
			},
			want: []Field{
				{
					Value: "",
					Name:  "$form_email",
					Data: map[string]interface{}{
						"label": "Email",
						"type":  "text",
					},
				},
			},
			wantErr: false,
		},

		{
			name: "TestForm_Parse: select type with option array #1",
			fields: fields{
				prefix: "_",
				code: `
$_lang = [
	'label' => 'Languages',
	'value' => 'php',
	'type' => 'select',
    'options' => [
        'c++' => 'cplusplus',
        'PHP' => 'php',
        'Go' => 'golang',
    ],
];
`,
			},
			want: []Field{
				{
					Value: "php",
					Name:  "$_lang",
					Data: map[string]interface{}{
						"label": "Languages",
						"type":  "select",
						"options": map[string]interface{}{
							"c++": "cplusplus",
							"PHP": "php",
							"Go":  "golang",
						},
					},
				},
			},
			wantErr: false,
		},

		{
			name: "TestForm_Parse: select type with option array #2",
			fields: fields{
				prefix: "_",
				code: `
$_lang = [
	'label' => 'Languages',
	'value' => 'php',
	'type' => 'select',
    'options' => [
        'cplusplus',
        'php',
        'golang',
    ],
];
`,
			},
			want: []Field{
				{
					Value: "php",
					Name:  "$_lang",
					Data: map[string]interface{}{
						"label": "Languages",
						"type":  "select",
						"options": []interface{}{
							"cplusplus",
							"php",
							"golang",
						},
					},
				},
			},
			wantErr: false,
		},

		{
			name: "TestForm_Parse: select type with option array #3",
			fields: fields{
				prefix: "_",
				code: `
$_lang = [
	'label' => 'Languages',
	'value' => 'php',
	'type' => 'select',
    'options' => [
		['label' => 'c++', 'value' => 'cplusplus'],
		['label' => 'PHP', 'value' => 'php'],
		['label' => 'Go', 'value' => 'golang'],
    ],
];
`,
			},
			want: []Field{
				{
					Value: "php",
					Name:  "$_lang",
					Data: map[string]interface{}{
						"label": "Languages",
						"type":  "select",
						"options": []interface{}{
							map[string]interface{}{
								"label": "c++",
								"value": "cplusplus",
							},
							map[string]interface{}{
								"label": "PHP",
								"value": "php",
							},
							map[string]interface{}{
								"label": "Go",
								"value": "golang",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "TestForm_Parse: option with function",
			fields: fields{
				prefix: "_",
				code: `
$_lang = [
	'label' => 'Languages',
	'value' => 'php',
	'type' => 'select',
    'options' => function() {
		return User::selectRaw('first_name as label, id as value')->get()->toArray();
	},
];
`,
			},
			want: []Field{
				{
					Value: "php",
					Name:  "$_lang",
					Data: map[string]interface{}{
						"label": "Languages",
						"type":  "select",
						"options": `function() {
		return User::selectRaw('first_name as label, id as value')->get()->toArray();
	}`,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "TestForm_Parse: option with fn",
			fields: fields{
				prefix: "_",
				code: `
$_lang = [
	'label' => 'Languages',
	'value' => 'php',
	'type' => 'select',
    'options' => fn() => Language::all(),
];
`,
			},
			want: []Field{
				{
					Value: "php",
					Name:  "$_lang",
					Data: map[string]interface{}{
						"label": "Languages",
						"type":  "select",
						"options": "fn() => Language::all()",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Form{
				prefix: tt.fields.prefix,
				code:   tt.fields.code,
			}
			got, err := f.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestForm_Stringify(t *testing.T) {
	type fields struct {
		prefix string
		code   string
	}
	type args struct {
		fields []Field
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "TestForm_Stringify",
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
				fields: []Field{
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
			want: `
$_name = [
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
		{
			name: "TestForm_Stringify: short",
			fields: fields{
				prefix: "_",
				code: `
$_name = 'billyct';
$_isAdmin = true;
$_age = 20;
`,
			},
			args: args{
				fields: []Field{
					{
						Name:  "$_name",
						Value: "magic",
					},
					{
						Name:  "$_isAdmin",
						Value: "false",
					},
					{
						Value: "30",
						Name:  "$_age",
					},
				},
			},
			want: `
$_name = 'magic';
$_isAdmin = false;
$_age = 30;`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Form{
				prefix: tt.fields.prefix,
				code:   tt.fields.code,
			}
			got, err := f.Stringify(tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stringify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Stringify() got = %v, want %v", got, tt.want)
			}
		})
	}
}
