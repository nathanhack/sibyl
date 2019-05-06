package core

import (
	"testing"
	"time"
)

func TestTimestampType_Before(t *testing.T) {

	b := TimestampType{0}
	a := TimestampType{1}

	if !b.Before(a) {
		t.Errorf("expected %v to be before %v", b, a)
	}

	b = NewTimestampTypeFromTime(time.Now().Add(-1 * time.Second))
	a = NewTimestampTypeFromTime(time.Now())

	if !b.Before(a) {
		t.Errorf("expected %v to be before %v", b, a)
	}
}

func TestTimestampType_After(t *testing.T) {

	b := TimestampType{0}
	a := TimestampType{1}

	if b.After(a) {
		t.Errorf("expected %v to be after %v", a, b)
	}

	b = NewTimestampTypeFromTime(time.Now().Add(-1 * time.Second))
	a = NewTimestampTypeFromTime(time.Now())

	if b.After(a) {
		t.Errorf("expected %v to be after %v", a, b)
	}
}

func TestTimestampType_Add(t *testing.T) {
	a := NewTimestampTypeFromTime(time.Now().Truncate(1 * time.Second))

	//since timestamp is only seconds based adding a nanosecond should do nothing
	b := a.Add(1 * time.Nanosecond)
	if !a.Equal(b) {
		t.Errorf("expected %v to be equal to %v", a, b)
	}

	b = a.Add(1 * time.Microsecond)
	if !a.Equal(b) {
		t.Errorf("expected %v to be equal to %v", a, b)
	}

	b = a.Add(1 * time.Millisecond)
	if !a.Equal(b) {
		t.Errorf("expected %v to be equal to %v", a, b)
	}

	c := a.Add(1 * time.Second)
	if !a.Before(c) {
		t.Errorf("expected %v to be before %v", a, c)
	}
}
