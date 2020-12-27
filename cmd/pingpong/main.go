package main

import (
	"fmt"
	"log"
	"time"

	"github.com/narslan/uci"
	"github.com/notnil/chess"
)

var engW *uci.Engine
var engB *uci.Engine

var resultOpts = uci.HighestDepthOnly | uci.IncludeUpperbounds | uci.IncludeLowerbounds

func init() {
	var err error
	engW, err = uci.NewEngine("./engines/stockfish")
	if err != nil {
		log.Fatal(err)
	}

	engB, err = uci.NewEngine("./engines/ethereal")
	if err != nil {
		log.Fatal(err)
	}

	// set some engine options
	engW.SetOptions(uci.Options{
		Hash:    128,
		Ponder:  true,
		OwnBook: true,
		MultiPV: 1,
		Threads: 4,
	})

	engB.SetOptions(uci.Options{
		Hash:    128,
		Ponder:  true,
		OwnBook: true,
		MultiPV: 1,
		Threads: 4,
	})
}

// The pinger prints a ping and waits for a pong
func wplay(white <-chan *chess.Game, black chan<- *chess.Game) {
	for {
		g := <-white

		bm := bestMoveWhite(g.FEN())
		fmt.Printf("[Stockfish]: %s\n", bm)

		if err := g.MoveStr(bm); err != nil {
			log.Printf("illegal move: %s\n", err.Error())
		}
		fmt.Println(g.Position().Board().Draw())
		time.Sleep(time.Second)
		black <- g
	}
}

// The ponger prints a pong and waits for a ping
func bplay(white chan<- *chess.Game, black <-chan *chess.Game) {
	for {
		g := <-black

		bm := bestMoveBlack(g.FEN())
		fmt.Printf("[Ethereal]: %s\n", bm)

		if err := g.MoveStr(bm); err != nil {
			log.Printf("illegal move: %s\n", err.Error())
		}
		fmt.Println(g.Position().Board().Draw())
		time.Sleep(time.Second)
		white <- g
	}
}

func main() {
	w := make(chan *chess.Game)
	b := make(chan *chess.Game)

	go wplay(w, b)
	go bplay(w, b)
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	g := startGame(fen)
	// The main goroutine starts the ping/pong by sending into the ping channel
	w <- g
	for {
		// Block the main thread until an interrupt
		time.Sleep(time.Second)
	}
}

func startGame(fens string) *chess.Game {

	fen, err := chess.FEN(fens)
	if err != nil {

		log.Printf("illegal fen: %s\n", err.Error())
	}

	return chess.NewGame(fen, chess.UseNotation(chess.LongAlgebraicNotation{}))
}

func bestMoveWhite(fens string) string {
	engW.SetFEN(fens)

	results, err := engW.GoDepth(10, resultOpts)
	if err != nil {
		fmt.Println(err.Error())
	}
	return results.BestMove
}

func bestMoveBlack(fens string) string {
	engB.SetFEN(fens)

	results, err := engB.GoDepth(10, resultOpts)
	if err != nil {
		fmt.Println(err.Error())
	}
	return results.BestMove
}
