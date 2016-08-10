package ctx

type TimeLimit struct {
	CurrentTimeMillis uint64
	DeadlineMillis    uint64
}

func (this *TimeLimit) SetCurrentTimeMillis(t uint64) (shouldAccept bool) {
	this.CurrentTimeMillis = t
	if t > this.DeadlineMillis {
		shouldAccept = false
	} else {
		shouldAccept = true
	}
	return
}

func (this TimeLimit) CalculateRestTimeMillis() uint64 {
	t := this.DeadlineMillis - this.CurrentTimeMillis
	if t < 0 {
		t = 0
	}
	return t
}
