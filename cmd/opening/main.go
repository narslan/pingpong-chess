package main

import (
	"fmt"

	"github.com/notnil/chess"
	"github.com/notnil/opening"
)

func main() {
	g := chess.NewGame()
	moves := []string{"e4", "c5", "b4"}
	for _, m := range moves {
		g.MoveStr(m)
	}
	o := opening.Find(g.Moves())
	fmt.Printf("%+v", o) // French Defense
}
