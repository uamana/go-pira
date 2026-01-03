package pira

import (
	"encoding/json"
	"testing"
)

func TestNullable_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		nullable Nullable[int]
		want     bool
	}{
		{
			name:     "valid nullable",
			nullable: Nullable[int]{Value: 42, Valid: true},
			want:     true,
		},
		{
			name:     "invalid nullable",
			nullable: Nullable[int]{Value: 0, Valid: false},
			want:     false,
		},
		{
			name:     "valid nullable with zero value",
			nullable: Nullable[int]{Value: 0, Valid: true},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.nullable.IsValid(); got != tt.want {
				t.Errorf("Nullable.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullable_IsNull(t *testing.T) {
	tests := []struct {
		name     string
		nullable Nullable[string]
		want     bool
	}{
		{
			name:     "null nullable",
			nullable: Nullable[string]{Value: "", Valid: false},
			want:     true,
		},
		{
			name:     "not null nullable",
			nullable: Nullable[string]{Value: "test", Valid: true},
			want:     false,
		},
		{
			name:     "not null nullable with empty string",
			nullable: Nullable[string]{Value: "", Valid: true},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.nullable.IsNull(); got != tt.want {
				t.Errorf("Nullable.IsNull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullable_String_Float64(t *testing.T) {
	tests := []struct {
		name     string
		nullable Nullable[float64]
		want     string
	}{
		{
			name:     "valid nullable with value",
			nullable: Nullable[float64]{Value: 123.45, Valid: true},
			want:     "123.45",
		},
		{
			name:     "invalid nullable",
			nullable: Nullable[float64]{Value: 0, Valid: false},
			want:     "not set",
		},
		{
			name:     "valid nullable with zero value",
			nullable: Nullable[float64]{Value: 0, Valid: true},
			want:     "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.nullable.String(); got != tt.want {
				t.Errorf("Nullable.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullable_String_String(t *testing.T) {
	tests := []struct {
		name     string
		nullable Nullable[string]
		want     string
	}{
		{
			name:     "valid nullable with string value",
			nullable: Nullable[string]{Value: "hello", Valid: true},
			want:     "hello",
		},
		{
			name:     "invalid nullable string",
			nullable: Nullable[string]{Value: "", Valid: false},
			want:     "not set",
		},
		{
			name:     "valid nullable with empty string",
			nullable: Nullable[string]{Value: "", Valid: true},
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.nullable.String(); got != tt.want {
				t.Errorf("Nullable.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullable_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		nullable Nullable[int]
		want     string
		wantErr  bool
	}{
		{
			name:     "valid nullable",
			nullable: Nullable[int]{Value: 42, Valid: true},
			want:     "42",
			wantErr:  false,
		},
		{
			name:     "invalid nullable",
			nullable: Nullable[int]{Value: 0, Valid: false},
			want:     "null",
			wantErr:  false,
		},
		{
			name:     "valid nullable with zero value",
			nullable: Nullable[int]{Value: 0, Valid: true},
			want:     "0",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nullable.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Nullable.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Nullable.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestNullable_MarshalJSON_Float64(t *testing.T) {
	tests := []struct {
		name     string
		nullable Nullable[float64]
		want     string
		wantErr  bool
	}{
		{
			name:     "valid nullable float",
			nullable: Nullable[float64]{Value: 123.45, Valid: true},
			want:     "123.45",
			wantErr:  false,
		},
		{
			name:     "invalid nullable float",
			nullable: Nullable[float64]{Value: 0, Valid: false},
			want:     "null",
			wantErr:  false,
		},
		{
			name:     "valid nullable with negative float",
			nullable: Nullable[float64]{Value: -123.45, Valid: true},
			want:     "-123.45",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nullable.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Nullable.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Nullable.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestNullable_MarshalJSON_String(t *testing.T) {
	tests := []struct {
		name     string
		nullable Nullable[string]
		want     string
		wantErr  bool
	}{
		{
			name:     "valid nullable string",
			nullable: Nullable[string]{Value: "hello", Valid: true},
			want:     `"hello"`,
			wantErr:  false,
		},
		{
			name:     "invalid nullable string",
			nullable: Nullable[string]{Value: "", Valid: false},
			want:     "null",
			wantErr:  false,
		},
		{
			name:     "valid nullable with empty string",
			nullable: Nullable[string]{Value: "", Valid: true},
			want:     `""`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nullable.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Nullable.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Nullable.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestNullable_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    Nullable[int]
		wantErr bool
	}{
		{
			name:    "valid integer",
			data:    []byte("42"),
			want:    Nullable[int]{Value: 42, Valid: true},
			wantErr: false,
		},
		{
			name:    "null value",
			data:    []byte("null"),
			want:    Nullable[int]{Value: 0, Valid: false},
			wantErr: false,
		},
		{
			name:    "zero value",
			data:    []byte("0"),
			want:    Nullable[int]{Value: 0, Valid: true},
			wantErr: false,
		},
		{
			name:    "negative value",
			data:    []byte("-42"),
			want:    Nullable[int]{Value: -42, Valid: true},
			wantErr: false,
		},
		{
			name:    "invalid json",
			data:    []byte("not a number"),
			want:    Nullable[int]{Value: 0, Valid: false},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Nullable[int]
			err := got.UnmarshalJSON(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Nullable.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Value != tt.want.Value || got.Valid != tt.want.Valid {
				t.Errorf("Nullable.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullable_UnmarshalJSON_Float64(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    Nullable[float64]
		wantErr bool
	}{
		{
			name:    "valid float",
			data:    []byte("123.45"),
			want:    Nullable[float64]{Value: 123.45, Valid: true},
			wantErr: false,
		},
		{
			name:    "null value",
			data:    []byte("null"),
			want:    Nullable[float64]{Value: 0, Valid: false},
			wantErr: false,
		},
		{
			name:    "zero value",
			data:    []byte("0"),
			want:    Nullable[float64]{Value: 0, Valid: true},
			wantErr: false,
		},
		{
			name:    "negative float",
			data:    []byte("-123.45"),
			want:    Nullable[float64]{Value: -123.45, Valid: true},
			wantErr: false,
		},
		{
			name:    "invalid json",
			data:    []byte("not a number"),
			want:    Nullable[float64]{Value: 0, Valid: false},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Nullable[float64]
			err := got.UnmarshalJSON(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Nullable.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Value != tt.want.Value || got.Valid != tt.want.Valid {
				t.Errorf("Nullable.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullable_UnmarshalJSON_String(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    Nullable[string]
		wantErr bool
	}{
		{
			name:    "valid string",
			data:    []byte(`"hello"`),
			want:    Nullable[string]{Value: "hello", Valid: true},
			wantErr: false,
		},
		{
			name:    "null value",
			data:    []byte("null"),
			want:    Nullable[string]{Value: "", Valid: false},
			wantErr: false,
		},
		{
			name:    "empty string",
			data:    []byte(`""`),
			want:    Nullable[string]{Value: "", Valid: true},
			wantErr: false,
		},
		{
			name:    "string with special chars",
			data:    []byte(`"hello@world#123"`),
			want:    Nullable[string]{Value: "hello@world#123", Valid: true},
			wantErr: false,
		},
		{
			name:    "invalid json",
			data:    []byte("not a string"),
			want:    Nullable[string]{Value: "", Valid: false},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Nullable[string]
			err := got.UnmarshalJSON(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Nullable.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Value != tt.want.Value || got.Valid != tt.want.Valid {
				t.Errorf("Nullable.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullable_UnmarshalJSON_Bool(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    Nullable[bool]
		wantErr bool
	}{
		{
			name:    "valid true",
			data:    []byte("true"),
			want:    Nullable[bool]{Value: true, Valid: true},
			wantErr: false,
		},
		{
			name:    "valid false",
			data:    []byte("false"),
			want:    Nullable[bool]{Value: false, Valid: true},
			wantErr: false,
		},
		{
			name:    "null value",
			data:    []byte("null"),
			want:    Nullable[bool]{Value: false, Valid: false},
			wantErr: false,
		},
		{
			name:    "invalid json",
			data:    []byte("not a boolean"),
			want:    Nullable[bool]{Value: false, Valid: false},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Nullable[bool]
			err := got.UnmarshalJSON(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Nullable.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Value != tt.want.Value || got.Valid != tt.want.Valid {
				t.Errorf("Nullable.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullable_JSON_RoundTrip(t *testing.T) {
	tests := []struct {
		name string
		val  Nullable[float64]
	}{
		{
			name: "valid value",
			val:  Nullable[float64]{Value: 123.45, Valid: true},
		},
		{
			name: "null value",
			val:  Nullable[float64]{Value: 0, Valid: false},
		},
		{
			name: "zero value but valid",
			val:  Nullable[float64]{Value: 0, Valid: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			data, err := json.Marshal(tt.val)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			// Unmarshal
			var got Nullable[float64]
			err = json.Unmarshal(data, &got)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			// Compare
			if got.Value != tt.val.Value || got.Valid != tt.val.Valid {
				t.Errorf("RoundTrip: got %v, want %v", got, tt.val)
			}
		})
	}
}

func TestNullable_JSON_RoundTrip_String(t *testing.T) {
	tests := []struct {
		name string
		val  Nullable[string]
	}{
		{
			name: "valid string",
			val:  Nullable[string]{Value: "hello world", Valid: true},
		},
		{
			name: "null value",
			val:  Nullable[string]{Value: "", Valid: false},
		},
		{
			name: "empty string but valid",
			val:  Nullable[string]{Value: "", Valid: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			data, err := json.Marshal(tt.val)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			// Unmarshal
			var got Nullable[string]
			err = json.Unmarshal(data, &got)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			// Compare
			if got.Value != tt.val.Value || got.Valid != tt.val.Valid {
				t.Errorf("RoundTrip: got %v, want %v", got, tt.val)
			}
		})
	}
}

func TestNullable_JSON_Struct(t *testing.T) {
	type TestStruct struct {
		ID    int
		Name  Nullable[string]
		Value Nullable[float64]
		Flag  Nullable[bool]
	}

	tests := []struct {
		name string
		data string
		want TestStruct
	}{
		{
			name: "all valid values",
			data: `{"ID":1,"Name":"test","Value":123.45,"Flag":true}`,
			want: TestStruct{
				ID:    1,
				Name:  Nullable[string]{Value: "test", Valid: true},
				Value: Nullable[float64]{Value: 123.45, Valid: true},
				Flag:  Nullable[bool]{Value: true, Valid: true},
			},
		},
		{
			name: "all null values",
			data: `{"ID":1,"Name":null,"Value":null,"Flag":null}`,
			want: TestStruct{
				ID:    1,
				Name:  Nullable[string]{Value: "", Valid: false},
				Value: Nullable[float64]{Value: 0, Valid: false},
				Flag:  Nullable[bool]{Value: false, Valid: false},
			},
		},
		{
			name: "mixed valid and null",
			data: `{"ID":1,"Name":"test","Value":null,"Flag":true}`,
			want: TestStruct{
				ID:    1,
				Name:  Nullable[string]{Value: "test", Valid: true},
				Value: Nullable[float64]{Value: 0, Valid: false},
				Flag:  Nullable[bool]{Value: true, Valid: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TestStruct
			err := json.Unmarshal([]byte(tt.data), &got)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			if got.ID != tt.want.ID {
				t.Errorf("ID = %v, want %v", got.ID, tt.want.ID)
			}
			if got.Name.Value != tt.want.Name.Value || got.Name.Valid != tt.want.Name.Valid {
				t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Value.Value != tt.want.Value.Value || got.Value.Valid != tt.want.Value.Valid {
				t.Errorf("Value = %v, want %v", got.Value, tt.want.Value)
			}
			if got.Flag.Value != tt.want.Flag.Value || got.Flag.Valid != tt.want.Flag.Valid {
				t.Errorf("Flag = %v, want %v", got.Flag, tt.want.Flag)
			}
		})
	}
}

