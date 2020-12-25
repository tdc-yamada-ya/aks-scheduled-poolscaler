package scaler

import (
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

func TestFindParameters(t *testing.T) {
	rules := Rules{
		{"0 0 1 1 2021 5", "p1"},
		{"* 18-23 * * * *", "p2"},
		{"* * * * * Sun,Sat", "p3"},
		{"* * * * * *", "p4"},
	}
	p1 := &Parameters{}
	p2 := &Parameters{}
	p3 := &Parameters{}
	p4 := &Parameters{}
	m := ParametersMap{
		"p1": p1,
		"p2": p2,
		"p3": p3,
		"p4": p4,
	}

	tests := []struct {
		t  time.Time
		p  *Parameters
		ok bool
	}{
		{time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC), p1, true},
		{time.Date(2021, time.January, 1, 18, 0, 0, 0, time.UTC), p2, true},
		{time.Date(2021, time.January, 2, 0, 0, 0, 0, time.UTC), p3, true},
		{time.Date(2021, time.January, 4, 0, 0, 0, 0, time.UTC), p4, true},
	}

	for _, tt := range tests {
		p, ok := FindParameters(rules, tt.t, m)

		if ok != tt.ok {
			t.Errorf("got %v, want %v", ok, tt.ok)
		}

		if p != tt.p {
			t.Errorf("got %v, want %v", p, tt.p)
		}
	}
}

var bpf = func(b bool) *bool { return &b }
var ipf = func(i int32) *int32 { return &i }

func TestUnmarshalConfigurationWithYaml(t *testing.T) {
	tests := []struct {
		f   string
		c   *Configuration
		err error
	}{
		{
			"testdata/configuration.yaml",
			&Configuration{
				ParametersMap{
					"p1": &Parameters{bpf(false), ipf(0), ipf(0), ipf(0)},
					"p2": &Parameters{EnableAutoScaling: bpf(true)},
					"p3": &Parameters{Count: ipf(1)},
					"p4": &Parameters{MinCount: ipf(1)},
					"p5": &Parameters{MaxCount: ipf(1)},
				},
				Resources{
					{"rgn1", "rn1", "apn1", Rules{{"0 0 1 Jan 2021 Fri", "p1"}, {"* * * * * *", "p2"}}},
					{"rgn2", "rn2", "apn2", Rules{{"0 0 1 Jan 2021 Fri", "p1"}, {"* * * * * *", "p2"}}},
				},
			},
			nil,
		},
	}

	for _, tt := range tests {
		b, err := ioutil.ReadFile(tt.f)
		if err != nil {
			t.Fatalf("read test configuration yaml error: %v", err)
		}

		c, err := UnmarshalConfigurationWithYaml(b)
		if !reflect.DeepEqual(c, tt.c) {
			t.Errorf("got %v, want %v", c, tt.c)
		}

		if err != tt.err {
			t.Errorf("got %v, want %v", err, tt.err)
		}
	}
}

func TestParametersString(t *testing.T) {
	tests := []struct {
		p *Parameters
		s string
	}{
		{&Parameters{nil, nil, nil, nil}, "EnableAutoScaling: nil, Count: nil, MinCount: nil, MaxCount: nil"},
		{&Parameters{bpf(true), ipf(1), ipf(1), ipf(1)}, "EnableAutoScaling: true, Count: 1, MinCount: 1, MaxCount: 1"},
	}

	for _, tt := range tests {
		s := tt.p.String()
		if s != tt.s {
			t.Errorf("got %v, want %v", s, tt.s)
		}
	}
}
