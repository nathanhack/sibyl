package internal

import (
	"time"
)

type TimeInterval struct {
	Start, End time.Time
}

// intervalDifference: will return intervals containing a-b meaning the intervals left after b interval has been removed
func IntervalDifference(a, b TimeInterval) []TimeInterval {

	// b doesn't overlap
	if !IntervalOverlap(a, b) {
		return []TimeInterval{a}
	}

	// b overlaps completely
	if (b.Start.Before(a.Start) || b.Start.Equal(a.Start)) &&
		(a.End.Before(b.End) || a.End.Equal(b.End)) {
		return []TimeInterval{}
	}

	// b overlaps at the beginning of a
	if b.Start.Before(a.Start) || b.Start.Equal(a.Start) {
		return []TimeInterval{{
			Start: b.End,
			End:   a.End,
		}}
	}

	// b overlaps at the end of a
	if a.End.Before(b.End) || a.End.Equal(b.End) {
		return []TimeInterval{{
			Start: a.Start,
			End:   b.Start,
		}}
	}

	// b must be completely inside
	return []TimeInterval{
		{Start: a.Start, End: b.Start},
		{Start: b.End, End: a.End},
	}
}

// intervalDifferences differences the first element with all the other elements
func IntervalDifferences(a ...TimeInterval) []TimeInterval {
	if len(a) == 0 {
		return []TimeInterval{}
	}

	if len(a) == 1 {
		return a
	}

	result := []TimeInterval{a[0]}
	for _, d := range a[1:] {
		tmp := make([]TimeInterval, 0)

		for _, r := range result {
			tmp = append(tmp, IntervalDifference(r, d)...)
		}

		result = tmp
	}
	return result
}

func IntervalDifferenceSlice(a []TimeInterval, b TimeInterval) []TimeInterval {
	result := make([]TimeInterval, 0)

	for _, r := range a {
		result = append(result, IntervalDifference(r, b)...)
	}

	return result
}

// intervalIntersection: will return intervals containing
func IntervalIntersection(a, b TimeInterval) []TimeInterval {

	// b and a don't overlap
	if !IntervalOverlap(a, b) {
		return []TimeInterval{}
	}

	result := TimeInterval{}

	if b.Start.Before(a.Start) {
		result.Start = a.Start
	} else if a.Start.Before(b.Start) {
		result.Start = b.Start
	} else {
		//a.Start.Equal(b.Start)
		result.Start = b.Start
	}

	if b.End.After(a.End) {
		result.End = a.End
	} else if a.End.After(b.End) {
		result.End = b.End
	} else {
		//a.End.Equal(b.End)
		result.End = b.End
	}

	return []TimeInterval{result}
}

func IntervalOverlap(a, b TimeInterval) bool {
	if a.Start.After(b.End) || a.End.Before(b.Start) {
		return false
	}
	return true
}

func IntervalUnion(a, b TimeInterval) []TimeInterval {
	if !IntervalOverlap(a, b) {
		return []TimeInterval{a, b}
	}

	min, max := a.Start, a.End
	if min.After(b.Start) {
		min = b.Start
	}

	if max.Before(b.End) {
		max = b.End
	}
	return []TimeInterval{{min, max}}
}
