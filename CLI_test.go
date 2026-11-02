package poker

import (
	"os"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	t.Run("check wins count", func(t *testing.T) {
		in := os.Stdin
		playerStore := &SpyStore{}
		cli := &CLI{playerStore, in}

		cli.PlayPoker()
		cli.PlayPoker()

		if len(playerStore.winCalls) != 2 {
			t.Fatal("expected a win call but didn't get any")
		}
	})

	t.Run("check wins count", func(t *testing.T) {
		in := strings.NewReader("Anna wins")
		playerStore := &SpyStore{}
		cli := &CLI{playerStore, in}

		cli.PlayPoker()
		got := playerStore.winCalls[0]
		want := "Anna"

		if got != want {
			t.Fatalf("expected a win call %q but got %q", want, got)
		}
	})
}
