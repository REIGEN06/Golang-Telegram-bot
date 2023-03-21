package main

import (
	"console_game/functions"
	"fmt"
)

func main() {
	fmt.Print("\x1b[2J")
	creature := functions.NewPlayer()
	fmt.Print("\tThe New Game is started!\n")
	creature.Info()
	for {
		creature.Day()
		if creature.CheckWin() {
			fmt.Print("Congratulations, you passed the game!")
			break
		}
		if creature.CheckDefeat() {
			fmt.Print("You lose the game")
			break
		}

		creature.Night()
		if creature.CheckWin() {
			fmt.Print("Congratulations, you passed the game!")
			break
		}
		if creature.CheckDefeat() {
			fmt.Print("You lose the game")
			break
		}
	}

}
