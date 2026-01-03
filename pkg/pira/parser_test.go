package pira

import (
	"testing"
)

func TestParseFloat64(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    float64
		wantErr bool
	}{
		{
			name:    "valid positive float",
			data:    []byte("123.45 "),
			want:    123.45,
			wantErr: false,
		},
		{
			name:    "valid positive float with spaces",
			data:    []byte("   123.45 "),
			want:    123.45,
			wantErr: false,
		},

		{
			name:    "valid negative float",
			data:    []byte("-123.45 "),
			want:    -123.45,
			wantErr: false,
		},
		{
			name:    "valid zero float",
			data:    []byte("0 "),
			want:    0,
			wantErr: false,
		},
		{
			name:    "valid float with scientific notation",
			data:    []byte("1.23e+02 "),
			want:    123.0,
			wantErr: false,
		},
		{
			name:    "valid float without trailing space",
			data:    []byte("123.45"),
			want:    123.45,
			wantErr: false,
		},
		{
			name:    "valid float with multiple spaces",
			data:    []byte("123.45  extra"),
			want:    123.45,
			wantErr: false,
		},
		{
			name:    "invalid float format",
			data:    []byte("not_a_number "),
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty value",
			data:    []byte(" "),
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty data",
			data:    []byte(""),
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFloat64(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseNullableFloat64(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want Nullable[float64]
	}{
		{
			name: "valid float with spaces",
			data: []byte("   123.45 "),
			want: Nullable[float64]{Value: 123.45, Valid: true},
		},
		{
			name: "valid float",
			data: []byte("123.45 "),
			want: Nullable[float64]{Value: 123.45, Valid: true},
		},
		{
			name: "valid float without trailing space",
			data: []byte("123.45"),
			want: Nullable[float64]{Value: 123.45, Valid: true},
		},
		{
			name: "invalid data",
			data: []byte("not_a_number "),
			want: Nullable[float64]{Value: 0, Valid: false},
		},
		{
			name: "empty data",
			data: []byte(""),
			want: Nullable[float64]{Value: 0, Valid: false},
		},
		{
			name: "only spaces",
			data: []byte("   "),
			want: Nullable[float64]{Value: 0, Valid: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseNullableFloat64(tt.data)
			if got.Value != tt.want.Value || got.Valid != tt.want.Valid {
				t.Errorf("parseNullableFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    int
		wantErr bool
	}{
		{
			name:    "valid positive int with spaces",
			data:    []byte("   123 "),
			want:    123,
			wantErr: false,
		},
		{
			name:    "valid positive int",
			data:    []byte("123 "),
			want:    123,
			wantErr: false,
		},
		{
			name:    "valid negative int",
			data:    []byte("-123 "),
			want:    -123,
			wantErr: false,
		},
		{
			name:    "valid zero int",
			data:    []byte("0 "),
			want:    0,
			wantErr: false,
		},
		{
			name:    "valid int without trailing space",
			data:    []byte("123"),
			want:    123,
			wantErr: false,
		},
		{
			name:    "valid int with multiple spaces",
			data:    []byte("123  extra"),
			want:    123,
			wantErr: false,
		},
		{
			name:    "invalid int format",
			data:    []byte("not_a_number "),
			want:    0,
			wantErr: true,
		},
		{
			name:    "float value (should fail)",
			data:    []byte("123.45 "),
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty value",
			data:    []byte(" "),
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty data",
			data:    []byte(""),
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseInt(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseNullableInt(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want Nullable[int]
	}{
		{
			name: "valid int with spaces",
			data: []byte("   123 "),
			want: Nullable[int]{Value: 123, Valid: true},
		},
		{
			name: "valid int",
			data: []byte("123 "),
			want: Nullable[int]{Value: 123, Valid: true},
		},
		{
			name: "valid int without trailing space",
			data: []byte("123"),
			want: Nullable[int]{Value: 123, Valid: true},
		},
		{
			name: "valid negative int",
			data: []byte("-123 "),
			want: Nullable[int]{Value: -123, Valid: true},
		},
		{
			name: "invalid data",
			data: []byte("not_a_number "),
			want: Nullable[int]{Value: 0, Valid: false},
		},
		{
			name: "empty data",
			data: []byte(""),
			want: Nullable[int]{Value: 0, Valid: false},
		},
		{
			name: "only spaces",
			data: []byte("   "),
			want: Nullable[int]{Value: 0, Valid: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseNullableInt(tt.data)
			if got.Value != tt.want.Value || got.Valid != tt.want.Valid {
				t.Errorf("parseNullableInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseString(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    string
		wantErr bool
	}{
		{
			name:    "valid string",
			data:    []byte("hello world "),
			want:    "hello world",
			wantErr: false,
		},
		{
			name:    "valid string with spaces",
			data:    []byte("  hello world  "),
			want:    "hello world",
			wantErr: false,
		},
		{
			name:    "valid string without trailing space",
			data:    []byte("hello world"),
			want:    "hello world",
			wantErr: false,
		},
		{
			name:    "valid empty string",
			data:    []byte(" "),
			want:    "",
			wantErr: false,
		},
		{
			name:    "valid string with special chars",
			data:    []byte("hello@world#123 "),
			want:    "hello@world#123",
			wantErr: false,
		},
		{
			name:    "valid string with multiple spaces",
			data:    []byte("hello world  extra"),
			want:    "hello world  extra",
			wantErr: false,
		},
		{
			name:    "empty data",
			data:    []byte(""),
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseString(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseNullableString(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want Nullable[string]
	}{
		{
			name: "valid string",
			data: []byte("hello world "),
			want: Nullable[string]{Value: "hello world", Valid: true},
		},
		{
			name: "valid string with leading/trailing spaces",
			data: []byte("  hello world  "),
			want: Nullable[string]{Value: "hello world", Valid: true},
		},
		{
			name: "empty data",
			data: []byte(""),
			want: Nullable[string]{Value: "", Valid: true},
		},
		{
			name: "only spaces",
			data: []byte("   "),
			want: Nullable[string]{Value: "", Valid: true},
		},
		{
			name: "string with special characters",
			data: []byte("hello@world#123 "),
			want: Nullable[string]{Value: "hello@world#123", Valid: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseNullableString(tt.data)
			if got.Value != tt.want.Value || got.Valid != tt.want.Valid {
				t.Errorf("parseNullableString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    bool
		wantErr bool
	}{
		{
			name:    "valid true - non-zero",
			data:    []byte("1 "),
			want:    true,
			wantErr: false,
		},
		{
			name:    "valid true - negative",
			data:    []byte("-1 "),
			want:    true,
			wantErr: false,
		},
		{
			name:    "valid false - zero",
			data:    []byte("0 "),
			want:    false,
			wantErr: false,
		},
		{
			name:    "valid true - large number",
			data:    []byte("42 "),
			want:    true,
			wantErr: false,
		},
		{
			name:    "valid bool without trailing space",
			data:    []byte("1"),
			want:    true,
			wantErr: false,
		},
		{
			name:    "invalid int format",
			data:    []byte("not_a_number "),
			want:    false,
			wantErr: true,
		},
		{
			name:    "empty value",
			data:    []byte(" "),
			want:    false,
			wantErr: true,
		},
		{
			name:    "empty data",
			data:    []byte(""),
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBool(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseNullableBool(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want Nullable[bool]
	}{
		{
			name: "valid true",
			data: []byte("1 "),
			want: Nullable[bool]{Value: true, Valid: true},
		},
		{
			name: "valid false",
			data: []byte("0 "),
			want: Nullable[bool]{Value: false, Valid: true},
		},
		{
			name: "valid true without trailing space",
			data: []byte("1"),
			want: Nullable[bool]{Value: true, Valid: true},
		},
		{
			name: "valid true with spaces",
			data: []byte("   1 "),
			want: Nullable[bool]{Value: true, Valid: true},
		},
		{
			name: "valid true with negative number",
			data: []byte("-1 "),
			want: Nullable[bool]{Value: true, Valid: true},
		},
		{
			name: "invalid data",
			data: []byte("not_a_number "),
			want: Nullable[bool]{Value: false, Valid: false},
		},
		{
			name: "empty data",
			data: []byte(""),
			want: Nullable[bool]{Value: false, Valid: false},
		},
		{
			name: "only spaces",
			data: []byte("   "),
			want: Nullable[bool]{Value: false, Valid: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseNullableBool(tt.data)
			if got.Value != tt.want.Value || got.Valid != tt.want.Valid {
				t.Errorf("parseNullableBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseHistogramData(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    Histogram
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid histogram with single pair",
			data:    []byte("0; 10"),
			want:    Histogram{[]int64{0, 10}},
			wantErr: false,
		},
		{
			name:    "valid histogram with multiple pairs",
			data:    []byte("0; 10 1; 20 2; 30"),
			want:    Histogram{[]int64{0, 10}, []int64{1, 20}, []int64{2, 30}},
			wantErr: false,
		},
		{
			name:    "valid histogram with multiple pairs and multiple lines",
			data:    []byte("0; 10 1;\r\n20 2; 30"),
			want:    Histogram{[]int64{0, 10}, []int64{1, 20}, []int64{2, 30}},
			wantErr: false,
		},

		{
			name:    "valid histogram with negative values",
			data:    []byte("-1; -5 0; 10"),
			want:    Histogram{[]int64{-1, -5}, []int64{0, 10}},
			wantErr: false,
		},
		{
			name:    "valid histogram with large numbers",
			data:    []byte("100; 1000 200; 2000"),
			want:    Histogram{[]int64{100, 1000}, []int64{200, 2000}},
			wantErr: false,
		},
		{
			name:    "valid histogram with extra spaces",
			data:    []byte("0;  10  1;  20"),
			want:    Histogram{[]int64{0, 10}, []int64{1, 20}},
			wantErr: false,
		},
		{
			name:    "empty data",
			data:    []byte(""),
			want:    Histogram{},
			wantErr: false,
		},
		{
			name:    "incomplete pair - only one field",
			data:    []byte("0;"),
			want:    Histogram{},
			wantErr: false,
		},
		{
			name:    "invalid bin - missing semicolon",
			data:    []byte("0 10"),
			want:    nil,
			wantErr: true,
			errMsg:  "invalid bin string",
		},
		{
			name:    "invalid bin - not a number",
			data:    []byte("abc; 10"),
			want:    nil,
			wantErr: true,
			errMsg:  "invalid bin string",
		},
		{
			name:    "invalid value - not a number",
			data:    []byte("0; abc"),
			want:    nil,
			wantErr: true,
			errMsg:  "invalid value string",
		},
		{
			name:    "invalid bin - float number",
			data:    []byte("0.5; 10"),
			want:    nil,
			wantErr: true,
			errMsg:  "invalid bin string",
		},
		{
			name:    "invalid value - float number",
			data:    []byte("0; 10.5"),
			want:    nil,
			wantErr: true,
			errMsg:  "invalid value string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHistogramData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHistogramData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if tt.errMsg != "" && err != nil && err.Error()[:len(tt.errMsg)] != tt.errMsg {
					t.Errorf("ParseHistogramData() error = %v, want error message containing %v", err, tt.errMsg)
				}
			} else {
				if len(got) != len(tt.want) {
					t.Errorf("ParseHistogramData() length = %v, want %v", len(got), len(tt.want))
					return
				}
				for i := range got {
					if len(got[i]) != len(tt.want[i]) {
						t.Errorf("ParseHistogramData() pair[%d] length = %v, want %v", i, len(got[i]), len(tt.want[i]))
						continue
					}
					for j := range got[i] {
						if got[i][j] != tt.want[i][j] {
							t.Errorf("ParseHistogramData() pair[%d][%d] = %v, want %v", i, j, got[i][j], tt.want[i][j])
						}
					}
				}
			}
		})
	}
}
