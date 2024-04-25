package time

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer := NewTimer(100)
	tds := make([]*TimerData, 100)
	for i := 0; i < 100; i++ {
		tds[i] = timer.Add(time.Duration(i)*time.Second+5*time.Minute, nil)
	}
	printTimer(t, timer)
	for i := 0; i < 100; i++ {
		t.Logf("td: %s, %s, %d", tds[i].Key, tds[i].ExpireString(), tds[i].index)
		timer.Del(tds[i])
	}
	printTimer(t, timer)
}
func printTimer(t *testing.T, timer *Timer) {
	t.Logf("----------timers: %d ----------", len(timer.timers))
	for i := 0; i < len(timer.timers); i++ {
		t.Logf("timer: %s, %s, index: %d", timer.timers[i].Key, timer.timers[i].ExpireString(), timer.timers[i].index)
	}
	t.Logf("--------------------")
}
