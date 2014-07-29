package internal

import (
	"reflect"
	"testing"
)

func TestInspectKind(t *testing.T) {
	tests := []struct {
		in     string
		expect reflect.Kind
	}{
		{"true", reflect.Bool},
		{"false", reflect.Bool},
		{"-1", reflect.Int64},
		{"0", reflect.Int64},
		{"1000", reflect.Int64},
		{"100.0", reflect.Float64},
		{"0.10", reflect.Float64},
		{"100e10", reflect.Float64},
		{"10e-10", reflect.Float64},
	}

	for n, tt := range tests {
		actual := InspectKind(tt.in)
		if actual != tt.expect {
			t.Errorf("#%d want %v, got %v", n, tt.expect, actual)
		}
	}
}
