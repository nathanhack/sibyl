package barrequester

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/nathanhack/sibyl/ent"
	"github.com/nathanhack/sibyl/cmd/server/cmd/internal"
)

func Test_findUniqueIntervalToBarTimeRanges(t *testing.T) {
	tests := []struct {
		input []*ent.BarTimeRange
		want  map[internal.TimeInterval][]*ent.BarTimeRange
	}{
		{ // nothing
			input: []*ent.BarTimeRange{},
			want:  map[internal.TimeInterval][]*ent.BarTimeRange{},
		},

		{ //no overlap
			input: []*ent.BarTimeRange{
				{ID: 0, Start: time.Date(2000, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2001, 1, 1, 12, 00, 00, 0, time.Local)},
			},
			want: map[internal.TimeInterval][]*ent.BarTimeRange{
				{Start: time.Date(2000, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2001, 1, 1, 12, 00, 00, 0, time.Local)}: {{ID: 0, Start: time.Date(2000, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2001, 1, 1, 12, 00, 00, 0, time.Local)}},
			},
		},
		{ //no overlap
			input: []*ent.BarTimeRange{
				{ID: 0, Start: time.Date(2000, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2001, 1, 1, 12, 00, 00, 0, time.Local)},
				{ID: 0, Start: time.Date(2002, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2003, 1, 1, 12, 00, 00, 0, time.Local)},
			},
			want: map[internal.TimeInterval][]*ent.BarTimeRange{
				{Start: time.Date(2000, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2001, 1, 1, 12, 00, 00, 0, time.Local)}: {{ID: 0, Start: time.Date(2000, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2001, 1, 1, 12, 00, 00, 0, time.Local)}},
				{Start: time.Date(2002, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2003, 1, 1, 12, 00, 00, 0, time.Local)}: {{ID: 0, Start: time.Date(2002, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2003, 1, 1, 12, 00, 00, 0, time.Local)}},
			},
		},

		{ // 1 overlap
			input: []*ent.BarTimeRange{
				{Start: time.Date(2000, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2001, 1, 1, 12, 00, 00, 0, time.Local)},
				{Start: time.Date(2001, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2003, 1, 1, 12, 00, 00, 0, time.Local)},
			},
			want: map[internal.TimeInterval][]*ent.BarTimeRange{
				{Start: time.Date(2000, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2003, 1, 1, 12, 00, 00, 0, time.Local)}: {
					{Start: time.Date(2000, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2001, 1, 1, 12, 00, 00, 0, time.Local)},
					{Start: time.Date(2001, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2003, 1, 1, 12, 00, 00, 0, time.Local)},
				},
			},
		},
		{ // overlap, nonoverlap,overlap
			input: []*ent.BarTimeRange{
				{Start: time.Date(2000, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2001, 1, 1, 12, 00, 00, 0, time.Local)},
				{Start: time.Date(2001, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2002, 1, 1, 12, 00, 00, 0, time.Local)},
				{Start: time.Date(2003, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2004, 1, 1, 12, 00, 00, 0, time.Local)},
				{Start: time.Date(2005, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2006, 1, 1, 12, 00, 00, 0, time.Local)},
				{Start: time.Date(2006, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2007, 1, 1, 12, 00, 00, 0, time.Local)},
			},
			want: map[internal.TimeInterval][]*ent.BarTimeRange{
				{Start: time.Date(2000, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2002, 1, 1, 12, 00, 00, 0, time.Local)}: {
					{Start: time.Date(2000, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2001, 1, 1, 12, 00, 00, 0, time.Local)},
					{Start: time.Date(2001, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2002, 1, 1, 12, 00, 00, 0, time.Local)},
				},
				{Start: time.Date(2003, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2004, 1, 1, 12, 00, 00, 0, time.Local)}: {
					{Start: time.Date(2003, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2004, 1, 1, 12, 00, 00, 0, time.Local)},
				},
				{Start: time.Date(2005, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2007, 1, 1, 12, 00, 00, 0, time.Local)}: {
					{Start: time.Date(2005, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2006, 1, 1, 12, 00, 00, 0, time.Local)},
					{Start: time.Date(2006, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2007, 1, 1, 12, 00, 00, 0, time.Local)},
				},
			},
		},
		{ // overlap, nonoverlap, overlap, nonoverlap
			input: []*ent.BarTimeRange{
				{Start: time.Date(2000, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2001, 1, 1, 12, 00, 00, 0, time.Local)},
				{Start: time.Date(2001, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2002, 1, 1, 12, 00, 00, 0, time.Local)},
				{Start: time.Date(2003, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2004, 1, 1, 12, 00, 00, 0, time.Local)},
				{Start: time.Date(2005, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2006, 1, 1, 12, 00, 00, 0, time.Local)},
				{Start: time.Date(2006, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2007, 1, 1, 12, 00, 00, 0, time.Local)},
				{Start: time.Date(2008, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2009, 1, 1, 12, 00, 00, 0, time.Local)},
			},
			want: map[internal.TimeInterval][]*ent.BarTimeRange{
				{Start: time.Date(2000, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2002, 1, 1, 12, 00, 00, 0, time.Local)}: {
					{Start: time.Date(2000, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2001, 1, 1, 12, 00, 00, 0, time.Local)},
					{Start: time.Date(2001, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2002, 1, 1, 12, 00, 00, 0, time.Local)},
				},
				{Start: time.Date(2003, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2004, 1, 1, 12, 00, 00, 0, time.Local)}: {
					{Start: time.Date(2003, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2004, 1, 1, 12, 00, 00, 0, time.Local)},
				},
				{Start: time.Date(2005, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2007, 1, 1, 12, 00, 00, 0, time.Local)}: {
					{Start: time.Date(2005, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2006, 1, 1, 12, 00, 00, 0, time.Local)},
					{Start: time.Date(2006, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2007, 1, 1, 12, 00, 00, 0, time.Local)},
				},
				{Start: time.Date(2008, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2009, 1, 1, 12, 00, 00, 0, time.Local)}: {
					{Start: time.Date(2008, 1, 1, 12, 00, 00, 0, time.Local), End: time.Date(2009, 1, 1, 12, 00, 00, 0, time.Local)},
				},
			},
		},
	}
	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			if got := findUniqueIntervalToBarTimeRanges(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findUniqueIntervalToBarTimeRanges() = %v, want %v", got, tt.want)
			}
		})
	}
}
