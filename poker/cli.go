package poker

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type CLI struct {
	playerStore PlayerStore
	in          *bufio.Scanner
	out         io.Writer
	game        Game
}

func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}

const PlayerPrompt = "Please enter the number of players: "

const BadPlayerInputErrMsg = "Bad value received for number of players, please try again with a number"

const BadWinnerInputMsg = "invalid winner input, expect format of 'PlayerName wins'"

func (cli *CLI) PlayPoker() error {
	if _, err := fmt.Fprint(cli.out, PlayerPrompt); err != nil {
		return err
	}

	numberOfPlayers, err := strconv.Atoi(cli.readLine())

	if err != nil {
		if _, err1 := fmt.Fprint(cli.out, BadPlayerInputErrMsg); err1 != nil {
			return err1
		}
		return nil
	}

	cli.game.Start(numberOfPlayers)

	winnerInput := cli.readLine()
	winner, err := extractWinner(winnerInput)

	if err != nil {
		if _, err1 := fmt.Fprint(cli.out, BadWinnerInputMsg); err1 != nil {
			return err1
		}
		return nil
	}

	if err := cli.game.Finish(winner); err != nil {
		return err
	}
	return nil
}

func extractWinner(userInput string) (string, error) {
	if !strings.Contains(userInput, " wins") {
		return "", errors.New(BadWinnerInputMsg)
	}
	return strings.Replace(userInput, " wins", "", 1), nil
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
