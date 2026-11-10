package poker

import (
	"fmt"
	"os"
	"time"
)

type Alert struct {
	scheduledAt time.Duration
	amount      int
}

func (a Alert) String() string {
	return fmt.Sprintf("%d chips at %v", a.amount, a.scheduledAt)
}

type BlindAlerter interface {
	ScheduleAlertAt(scheduledAt time.Duration, amount int)
}

type BlindAlerterFunc func(duration time.Duration, amount int)

func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int) {
	a(duration, amount)
}

func StdOutAlerter(duration time.Duration, amount int) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(os.Stdout, "Blind is now %d\n", amount)
	})
}
