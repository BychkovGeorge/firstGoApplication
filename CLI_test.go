package poker

import (
	"os"
	"strings"
	"testing"
	"time"
)

type SpyBlindAlerter struct {
	alerts []Alert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, Alert{duration, amount})
}

func TestCLI(t *testing.T) {
	t.Run("check wins count", func(t *testing.T) {
		in := os.Stdin
		playerStore := &SpyStore{}
		dummyBlindAlerter := &SpyBlindAlerter{}
		cli := NewCLI(playerStore, in, dummyBlindAlerter)

		cli.PlayPoker()
		cli.PlayPoker()

		if len(playerStore.winCalls) != 2 {
			t.Fatal("expected a win call but didn't get any")
		}
	})

	t.Run("check wins count", func(t *testing.T) {
		in := strings.NewReader("Anna wins")
		playerStore := &SpyStore{}
		dummyBlindAlerter := &SpyBlindAlerter{}
		cli := NewCLI(playerStore, in, dummyBlindAlerter)

		cli.PlayPoker()
		got := playerStore.winCalls[0]
		want := "Anna"

		if got != want {
			t.Fatalf("expected a win call %q but got %q", want, got)
		}
	})

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &SpyStore{}
		blindAlerter := &SpyBlindAlerter{}

		cli := NewCLI(playerStore, in, blindAlerter)
		cli.PlayPoker()

		if len(blindAlerter.alerts) != 1 {
			t.Fatal("expected a blind alert to be scheduled")
		}
	})
}
