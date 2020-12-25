package main

import (
	"testing"
)

func TestOptionsValidate(t *testing.T) {
	tests := []struct {
		o options
		e bool
	}{
		{options{"", "", ""}, true},
		{options{"a", "", ""}, true},
		{options{"", "a", ""}, true},
		{options{"", "", "a"}, true},
		{options{"a", "a", "a"}, false},
	}

	for _, tt := range tests {
		err := tt.o.validate()
		e := err != nil
		if e != tt.e {
			t.Errorf("got %v, want %v", e, tt.e)
		}
	}
}
