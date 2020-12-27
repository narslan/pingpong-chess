package main

import (
	"fmt"
	"log"

	"github.com/narslan/uci"
	"github.com/notnil/chess"
)

var eng *uci.Engine
var resultOpts = uci.HighestDepthOnly | uci.IncludeUpperbounds | uci.IncludeLowerbounds

func init() {
	var err error
	eng, err = uci.NewEngine("./engines/stockfish")
	if err != nil {
		log.Fatal(err)
	}

	// set some engine options
	eng.SetOptions(uci.Options{
		Hash:    128,
		Ponder:  false,
		OwnBook: true,
		MultiPV: 1,
	})

}

func main() {

	fen := "rn1q1rk1/pbn1bppp/1pp5/3p1N2/8/2N3P1/PP2PPBP/R1BQ1RK1 w - d6 0 12"
	g := startGame(fen)

	bm := bestMove(g.FEN())
	fmt.Printf("[Stockfish]: %s", bm)

	if err := g.MoveStr(bm); err != nil {
		log.Printf("illegal move: %s\n", err.Error())
	}

	fmt.Println(g.Position().Board().Draw())

	// if err := g.MoveStr("Kf8"); err != nil {
	// 	log.Printf("illegal move: %s\n", err.Error())
	// }

	//fmt.Println(g.Position().Board().Draw())

}

func bestMove(fens string) string {
	eng.SetFEN(fens)

	results, err := eng.GoDepth(10, resultOpts)
	if err != nil {
		fmt.Println(err.Error())
	}
	return results.BestMove

}

func startGame(fens string) *chess.Game {

	fen, err := chess.FEN(fens)
	if err != nil {

		log.Printf("illegal fen: %s\n", err.Error())
	}

	return chess.NewGame(fen, chess.UseNotation(chess.LongAlgebraicNotation{}))
}
