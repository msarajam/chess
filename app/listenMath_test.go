package app

import (
	"reflect"
	"testing"

	"webServer/config"
)

func Test_checkFormula(t *testing.T) {
	var f config.MathFormulaStruct
	tests := []struct {
		name  string
		s     []string
		want1 bool
	}{
		{
			name:  "Successful Int value",
			s:     []string{"1", "2"},
			want1: true,
		},
		{
			name:  "Successful Float value",
			s:     []string{"1.0", "2.00"},
			want1: true,
		},
		{
			name:  "Fail Char First Value",
			s:     []string{"A", "1"},
			want1: false,
		},
		{
			name:  "Fail Char Second Value",
			s:     []string{"1", "A"},
			want1: false,
		},
		{
			name:  "Fail One Value",
			s:     []string{"1"},
			want1: false,
		},

		{
			name:  "Fail First Value Empty",
			s:     []string{"", "1"},
			want1: false,
		},
		{
			name:  "Fail Second Value empty",
			s:     []string{"1", ""},
			want1: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := checkFormula(tt.s, &f)
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("checkFormula got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func Test_findOperation(t *testing.T) {
	var f config.MathFormulaStruct
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "Successful Plus +",
			s:    "+",
			want: false,
		},

		{
			name: "Successful Negative -",
			s:    "-",
			want: false,
		},
		{
			name: "Successful Multiply *",
			s:    "*",
			want: false,
		},
		{
			name: "Successful Devide /",
			s:    "/",
			want: false,
		},
		{
			name: "Fail Empyt Value",
			s:    "",
			want: true,
		},
		{
			name: "Fail Wrong Value",
			s:    "1",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findOperation(&tt.s, &f)
			if got != nil && !tt.want {
				t.Errorf("findOperation got = %v, want: %v", got, tt.want)
			}
		})
	}
}

func Test_calculation(t *testing.T) {
	tests := []struct {
		name    string
		f       config.MathFormulaStruct
		want    string
		wantErr bool
	}{
		{
			name: "Successful Plus +",
			f: config.MathFormulaStruct{
				FirstValue:  1.0,
				SecondValue: 2.0,
				Operation:   "+",
				Answer:      "",
			},
			want:    "3.00",
			wantErr: false,
		},
		{
			name: "Successful Negative -",
			f: config.MathFormulaStruct{
				FirstValue:  2.0,
				SecondValue: 1.0,
				Operation:   "-",
				Answer:      "",
			},
			want:    "1.00",
			wantErr: false,
		},
		{
			name: "Successful Multiply *",
			f: config.MathFormulaStruct{
				FirstValue:  1.0,
				SecondValue: 2.0,
				Operation:   "*",
				Answer:      "",
			},
			want:    "2.00",
			wantErr: false,
		},
		{
			name: "Successful Devide /",
			f: config.MathFormulaStruct{
				FirstValue:  4.0,
				SecondValue: 2.0,
				Operation:   "/",
				Answer:      "",
			},
			want:    "2.00",
			wantErr: false,
		},
		{
			name: "Fail Plus +",
			f: config.MathFormulaStruct{
				FirstValue:  3.0,
				SecondValue: 2.0,
				Operation:   "+",
				Answer:      "",
			},
			want:    "3.00",
			wantErr: true,
		},
		{
			name: "Fail Negative -",
			f: config.MathFormulaStruct{
				FirstValue:  1.0,
				SecondValue: 2.0,
				Operation:   "-",
				Answer:      "",
			},
			want:    "3.00",
			wantErr: true,
		},
		{
			name: "Fail Multiply *",
			f: config.MathFormulaStruct{
				FirstValue:  1.0,
				SecondValue: 2.0,
				Operation:   "*",
				Answer:      "",
			},
			want:    "3.00",
			wantErr: true,
		},
		{
			name: "Fail Devide /",
			f: config.MathFormulaStruct{
				FirstValue:  1.0,
				SecondValue: 2.0,
				Operation:   "/",
				Answer:      "",
			},
			want:    "3.00",
			wantErr: true,
		},
		{
			name: "Fail Wrong Operation",
			f: config.MathFormulaStruct{
				FirstValue:  1.0,
				SecondValue: 2.0,
				Operation:   "A",
				Answer:      "",
			},
			want:    "3.00",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := calculation(&tt.f); err == nil && !tt.wantErr {
				if tt.f.Answer != tt.want && tt.wantErr {
					t.Errorf("calculation got = %v, want: %v", tt.f.Answer, tt.want)
				}
			}
		})
	}
}
