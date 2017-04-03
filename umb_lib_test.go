package umb_lib

import (
	"fmt"
	"testing"
	"time"
)

func TestReadConf(t *testing.T) {
	conf := ReadConf("")

	cases := []struct{ key, want string }{
		{"DBUSER", "tamaex"},
		{"QWE", ""},
		{"", ""},
	}

	for _, c := range cases {
		fmt.Println("key=", c.key, ", got=", conf[c.key], ", want=", c.want)
		if conf[c.key] != c.want {
			t.Errorf("Error on ReadConf: want=%s, got=%s\n", c.want, conf[c.key])
		}
	}
}

func TestTimeConv(t *testing.T) {
	cases := []struct {
		tt time.Time
		ts string
	}{
		{time.Date(2010, time.June, 13, 0, 0, 0, 0, time.UTC), "13.06.2010"},
	}

	for _, c := range cases {
		got := T2D(c.tt)
		fmt.Printf("tt=%#v, ts=%s, got=%s\n", c.tt, c.ts, got)
		if got != c.ts {
			t.Errorf("Error on T2D: want=%s, got=%s\n", c.ts, got)
		}
		tt2, err := D2T(c.ts)
		fmt.Printf("tt=%#v, tt2=%#v\n", c.tt, tt2)
		if tt2 != c.tt || err != nil {
			t.Errorf("Error on D2T: \twant=%#v, got=%#v\n\terr=%s\n", c.tt, tt2, err)
		}
	}
}

func TestRound(t *testing.T) {
	cases := []struct {
		in1  float64
		in2  int
		want float64
	}{
		{1.234, 2, 1.23},
		{1.236, 2, 1.24},
		{1.234, 0, 1},
		{0, 0, 0},
	}

	for _, c := range cases {
		got := Round(c.in1, c.in2)
		fmt.Printf("Error on Round(%f,%d)=%f <?> %f\n", c.in1, c.in2, got, c.want)
		if got != c.want {
			t.Errorf("Error on Round(%f,%d)=%f != %f\n", c.in1, c.in2, got, c.want)
		}
	}
}
