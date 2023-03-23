package internal

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

func Test_intervalIntersection(t *testing.T) {
	type args struct {
		a TimeInterval
		b TimeInterval
	}
	tests := []struct {
		args args
		want []TimeInterval
	}{

		0: {args: args{
			a: TimeInterval{Start: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local)},
			b: TimeInterval{Start: time.Date(2001, 1, 2, 0, 0, 0, 0, time.Local), End: time.Date(2009, 1, 1, 0, 0, 0, 0, time.Local)},
		},
			want: []TimeInterval{},
		},
		1: {args: args{
			a: TimeInterval{Start: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local)},
			b: TimeInterval{Start: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2009, 1, 1, 0, 0, 0, 0, time.Local)},
		},
			want: []TimeInterval{{Start: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local)}},
		},
		2: {args: args{
			a: TimeInterval{Start: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2010, 1, 1, 0, 0, 0, 0, time.Local)},
			b: TimeInterval{Start: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2009, 1, 1, 0, 0, 0, 0, time.Local)},
		},
			want: []TimeInterval{{Start: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2009, 1, 1, 0, 0, 0, 0, time.Local)}},
		},
		3: {args: args{
			a: TimeInterval{Start: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2010, 1, 1, 0, 0, 0, 0, time.Local)},
			b: TimeInterval{Start: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2009, 1, 1, 0, 0, 0, 0, time.Local)},
		},
			want: []TimeInterval{{Start: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2009, 1, 1, 0, 0, 0, 0, time.Local)}},
		},
		4: {args: args{
			a: TimeInterval{Start: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2009, 1, 1, 0, 0, 0, 0, time.Local)},
			b: TimeInterval{Start: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2010, 1, 1, 0, 0, 0, 0, time.Local)},
		},
			want: []TimeInterval{{Start: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2009, 1, 1, 0, 0, 0, 0, time.Local)}},
		},
	}
	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			if got := IntervalIntersection(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("intervalIntersection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_intervalDifference(t *testing.T) {
	type args struct {
		a TimeInterval
		b TimeInterval
	}
	tests := []struct {
		name string
		args args
		want []TimeInterval
	}{
		{args: args{
			a: TimeInterval{Start: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local)},
			b: TimeInterval{Start: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2009, 1, 1, 0, 0, 0, 0, time.Local)},
		},
			want: []TimeInterval{{Start: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local)}},
		},

		{args: args{
			a: TimeInterval{Start: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2010, 1, 1, 0, 0, 0, 0, time.Local)},
			b: TimeInterval{Start: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2010, 1, 1, 0, 0, 0, 0, time.Local)},
		},
			want: []TimeInterval{},
		},

		{args: args{
			a: TimeInterval{Start: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2010, 1, 1, 0, 0, 0, 0, time.Local)},
			b: TimeInterval{Start: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2009, 1, 1, 0, 0, 0, 0, time.Local)},
		},
			want: []TimeInterval{
				{Start: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local)},
				{Start: time.Date(2009, 1, 1, 0, 0, 0, 0, time.Local), End: time.Date(2010, 1, 1, 0, 0, 0, 0, time.Local)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntervalDifference(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("intervalDifference() = %v, want %v", got, tt.want)
			}
		})
	}
}
