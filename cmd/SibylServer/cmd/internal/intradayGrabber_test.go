package internal

import (
	"testing"
	"time"
)

func TestIntradayStockRange_ReduceDelta(t *testing.T) {
	isr := intradayStockRange{
		Delta:     defaultISRStartingDelta,
		StartDate: time.Now().Add(-defaultIntradyStart),
		EndDate:   time.Now(),
	}

	for isr.ReduceDelta() {
		t.Log(isr)
	}

	if isr.Delta != 0 {
		t.Errorf("expected zero but found %v", isr.Delta)
	}

}
