package reflection

import "testing"

func TestDescribeValue(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{name: "struct value", input: taggedUser{Name: "gopher", Level: 2}, want: "type=reflection.taggedUser kind=struct value={gopher 2}"},
		{name: "nil value", input: nil, want: "invalid value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := describeValue(tt.input); got != tt.want {
				t.Fatalf("describeValue() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReadStructTags(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{name: "struct tags", input: taggedUser{}, want: "Name json=name db=user_name | Level json=level db=user_level"},
		{name: "non struct", input: 42, want: "no tags"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readStructTags(tt.input); got != tt.want {
				t.Fatalf("readStructTags() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCallMethod(t *testing.T) {
	tests := []struct {
		name   string
		target any
		method string
		args   []string
		want   string
	}{
		{name: "dynamic greet", target: greeter{Prefix: "hi"}, method: "Greet", args: []string{"Go"}, want: "hi, Go"},
		{name: "missing method", target: greeter{Prefix: "hi"}, method: "Wave", want: "method not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := callMethod(tt.target, tt.method, tt.args...); got != tt.want {
				t.Fatalf("callMethod() = %q, want %q", got, tt.want)
			}
		})
	}
}
