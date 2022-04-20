package config

import (
	"reflect"
	"testing"
	"time"
)

func TestParseInterval(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		arg     string
		wantI   *Interval
		wantErr bool
	}{
		{
			"single integer",
			"YEAR=2020",
			&Interval{year: []int{2020}},
			false,
		},
		{
			"multiple integers",
			"MONTH=1,3,5",
			&Interval{month: []int{1, 3, 5}},
			false,
		},
		{
			"range of integers",
			"day=7-12",
			&Interval{day: []int{7, 8, 9, 10, 11, 12}},
			false,
		},
		{
			"out of sequence",
			"hour=6,4,2",
			&Interval{hour: []int{2, 4, 6}},
			false,
		},
		{
			"duplicates",
			"minute=30,40,30",
			&Interval{minute: []int{30, 40}},
			false,
		},
		{
			"weekdays",
			"weekday=f,weds,Thursday,w",
			&Interval{weekday: []int{3, 4, 5}},
			false,
		},
		{
			"multiple terms",
			"month=3,6 weekday=m",
			&Interval{month: []int{3, 6}, weekday: []int{1}},
			false,
		},
		{
			"empty string",
			"  ",
			nil, true,
		},
		{
			"number out of range",
			"MONTH=13",
			nil, true,
		},
		{
			"range out of order",
			"YEAR=2030-2020",
			nil, true,
		},
		{
			"unknown weekday",
			"WEEKDAY=x",
			nil, true,
		},
		{
			"ambiguous weekday",
			"WEEKDAY=t",
			nil, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotI, err := ParseInterval(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInterval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotI, tt.wantI) {
				t.Errorf("ParseInterval() = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}

func TestInterval_Next(t *testing.T) {
	tests := []struct {
		name     string
		interval string
		from     time.Time
		wantNext time.Time
	}{
		{
			"every 30 minutes",
			"MINUTE=0,30",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			time.Date(2020, 1, 1, 0, 30, 0, 0, time.Local),
		},
		{
			"crossing date boundary",
			"DAY=1,3,5 MINUTE=0,30",
			time.Date(2020, 1, 1, 23, 45, 0, 0, time.Local),
			time.Date(2020, 1, 3, 0, 0, 0, 0, time.Local),
		},
		{
			"crossing month boundary",
			"MONTH=1,3,5 MINUTE=0,30",
			time.Date(2020, 1, 31, 23, 45, 0, 0, time.Local),
			time.Date(2020, 3, 1, 0, 0, 0, 0, time.Local),
		},
		{
			"8:00pm on the second Tuesday of the month",
			"DAY=8-14 WEEKDAY=Tu HOUR=20 MINUTE=0",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			time.Date(2020, 1, 14, 20, 0, 0, 0, time.Local),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, err := ParseInterval(tt.interval)
			if err != nil {
				t.Errorf("parse failed")
			}
			if gotNext := i.Next(tt.from); !reflect.DeepEqual(gotNext, tt.wantNext) {
				t.Errorf("Interval.Next() = %v, want %v", gotNext, tt.wantNext)
			}
		})
	}
}

func TestInterval_Prev(t *testing.T) {
	tests := []struct {
		name     string
		interval string
		from     time.Time
		wantPrev time.Time
	}{
		{
			"every 30 minutes",
			"MINUTE=0,30",
			time.Date(2020, 1, 1, 0, 30, 0, 0, time.Local),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
		},
		{
			"crossing date boundary",
			"DAY=1,3,5 MINUTE=0,30",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			time.Date(2019, 12, 5, 23, 30, 0, 0, time.Local),
		},
		{
			"crossing month boundary",
			"MONTH=1,3,5 MINUTE=0,30",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			time.Date(2019, 5, 31, 23, 30, 0, 0, time.Local),
		},
		{
			"8:00pm on the second Tuesday of the month",
			"DAY=8-14 WEEKDAY=Tu HOUR=20 MINUTE=0",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			time.Date(2019, 12, 10, 20, 0, 0, 0, time.Local),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, err := ParseInterval(tt.interval)
			if err != nil {
				t.Errorf("parse failed")
			}
			if gotPrev := i.Prev(tt.from); !reflect.DeepEqual(gotPrev, tt.wantPrev) {
				t.Errorf("Interval.Prev() = %v, want %v", gotPrev, tt.wantPrev)
			}
		})
	}
}
