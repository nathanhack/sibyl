// Code generated by ogen, DO NOT EDIT.

package ogent

// setDefaults set default value of fields.
func (s *BarGroupTimeRangeRead) setDefaults() {
	{
		val := BarGroupTimeRangeReadStatus("pending")
		s.Status = val
	}
}

// setDefaults set default value of fields.
func (s *BarTimeRangeCreate) setDefaults() {
	{
		val := BarTimeRangeCreateStatus("pending")
		s.Status = val
	}
}

// setDefaults set default value of fields.
func (s *BarTimeRangeList) setDefaults() {
	{
		val := BarTimeRangeListStatus("pending")
		s.Status = val
	}
}

// setDefaults set default value of fields.
func (s *BarTimeRangeRead) setDefaults() {
	{
		val := BarTimeRangeReadStatus("pending")
		s.Status = val
	}
}

// setDefaults set default value of fields.
func (s *BarTimeRangeUpdate) setDefaults() {
	{
		val := BarTimeRangeUpdateStatus("pending")
		s.Status = val
	}
}

// setDefaults set default value of fields.
func (s *CreateBarTimeRangeReq) setDefaults() {
	{
		val := CreateBarTimeRangeReqStatus("pending")
		s.Status = val
	}
}

// setDefaults set default value of fields.
func (s *IntervalBarsList) setDefaults() {
	{
		val := IntervalBarsListStatus("pending")
		s.Status = val
	}
}

// setDefaults set default value of fields.
func (s *UpdateBarTimeRangeReq) setDefaults() {
	{
		val := UpdateBarTimeRangeReqStatus("pending")
		s.Status.SetTo(val)
	}
}
