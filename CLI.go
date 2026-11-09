package poker

import (
	"bufio"
	"io"
	"strings"
	"time"
)

type Alert struct {
	scheduledAt time.Duration
	amount      int
}

type BlindAlerter interface {
	ScheduleAlertAt(scheduledAt time.Duration, amount int)
}

type CLI struct {
	playerStore PlayerStore
	in          io.Reader
	alerter     BlindAlerter
}

func NewCLI(playerStore PlayerStore, in io.Reader, blindAlerter BlindAlerter) *CLI {
	return &CLI{playerStore, in, blindAlerter}
}

func (cli *CLI) PlayPoker() {
	cli.alerter.ScheduleAlertAt(5*time.Second, 100)
	reader := bufio.NewScanner(cli.in)
	reader.Scan()
	cli.playerStore.IncrementPlayerScore(extractWinner(reader.Text()))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}
