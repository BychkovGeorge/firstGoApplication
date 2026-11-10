package poker

import (
	"bytes"
	"fmt"
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
	var dummyBlindAlerter = &SpyBlindAlerter{}
	var dummyPlayerStore = &SpyStore{}
	var dummyStdIn = &bytes.Buffer{}
	var dummyStdOut = &bytes.Buffer{}

	t.Run("check wins count", func(t *testing.T) {
		in := os.Stdin
		playerStore := &SpyStore{}
		dummyBlindAlerter := &SpyBlindAlerter{}
		cli := NewCLI(in, dummyStdOut, NewGame(dummyBlindAlerter, playerStore))

		cli.PlayPoker()
		cli.PlayPoker()

		if len(playerStore.winCalls) != 2 {
			t.Fatal("expected a win call but didn't get any")
		}
	})

	t.Run("check wins count", func(t *testing.T) {
		in := strings.NewReader("5\nAnna wins\n")
		playerStore := &SpyStore{}
		dummyBlindAlerter := &SpyBlindAlerter{}
		cli := NewCLI(in, dummyStdOut, NewGame(dummyBlindAlerter, playerStore))

		cli.PlayPoker()
		got := playerStore.winCalls[0]
		want := "Anna"

		if got != want {
			t.Fatalf("expected a win call %q but got %q", want, got)
		}
	})

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("5\nChris wins\n")
		playerStore := &SpyStore{}
		blindAlerter := &SpyBlindAlerter{}

		cli := NewCLI(in, dummyStdOut, NewGame(blindAlerter, playerStore))
		cli.PlayPoker()

		cases := []Alert{
			{0 * time.Second, 100},
			{10 * time.Second, 200},
			{20 * time.Second, 300},
			{30 * time.Second, 400},
			{40 * time.Second, 500},
			{50 * time.Second, 600},
			{60 * time.Second, 800},
			{70 * time.Second, 1000},
			{80 * time.Second, 2000},
			{90 * time.Second, 4000},
			{100 * time.Second, 8000},
		}

		for i, want := range cases {
			t.Run(fmt.Sprintf("%d scheduled for %v", want.amount, want.scheduledAt), func(t *testing.T) {

				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduledAlert(t, got, want)
			})
		}
	})

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		cli := NewCLI(dummyStdIn, stdout, NewGame(dummyBlindAlerter, dummyPlayerStore))
		cli.PlayPoker()

		got := stdout.String()
		want := PlayerPrompt

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\nAnna wins\n")
		blindAlerter := &SpyBlindAlerter{}

		cli := NewCLI(in, stdout, NewGame(blindAlerter, dummyPlayerStore))
		cli.PlayPoker()

		got := stdout.String()
		want := PlayerPrompt

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		cases := []Alert{
			{0 * time.Second, 100},
			{12 * time.Second, 200},
			{24 * time.Second, 300},
			{36 * time.Second, 400},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {

				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduledAlert(t, got, want)
			})
		}
	})
}

func assertScheduledAlert(t *testing.T, got, want Alert) {
	amountGot := got.amount
	if amountGot != want.amount {
		t.Errorf("got amount %d, want %d", amountGot, want.amount)
	}

	gotScheduledTime := got.scheduledAt
	if gotScheduledTime != want.scheduledAt {
		t.Errorf("got scheduled time of %v, want %v", gotScheduledTime, want.scheduledAt)
	}
}
