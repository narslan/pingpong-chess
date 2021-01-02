package main

import (
	"fmt"
	"log"
	_ "time"

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

	engB, err = uci.NewEngine("./engines/ethereal_nnue")
	if err != nil {
		log.Fatal(err)
	}
	optW := map[string]interface{}{
		"MultiPV":  4,
		"Threads":  4,
		"EvalFile": "/home/nevroz/go/src/github.com/narslan/pingpong-chess/engines/nn-82215d0fd0df.nnue",
	}
	optB := map[string]interface{}{
		"MultiPV":  4,
		"Threads":  4,
		"EvalFile": "/home/nevroz/go/src/github.com/narslan/pingpong-chess/engines/nn-82215d0fd0df.nnue",
	}

	// set some engine options
	engW.SetOptions(optW)
	engB.SetOptions(optB)
}

// The pinger prints a ping and waits for a pong
func wplay(white <-chan *chess.Game, black chan<- *chess.Game, done chan bool) {
	for {
		g := <-white

		bm := bestMoveWhite(g.FEN())
		fmt.Printf("[Stockfish]: %s\n", bm)

		if err := g.MoveStr(bm); err != nil {
			log.Printf("illegal move: %s\n", err.Error())
			done <- true
		}
		fmt.Println(g.Position().Board().Draw())
		//time.Sleep(1 * time.Second)
		black <- g
	}
}

// The ponger prints a pong and waits for a ping
func bplay(white chan<- *chess.Game, black <-chan *chess.Game, done chan bool) {
	for {
		g := <-black

		bm := bestMoveBlack(g.FEN())
		fmt.Printf("[Ethereal]: %s\n", bm)

		if err := g.MoveStr(bm); err != nil {
			log.Printf("illegal move: %s\n", err.Error())
			done <- true
		}
		fmt.Println(g.Position().Board().Draw())
		//time.Sleep(1 * time.Second)
		white <- g
	}
}

func main() {
	w := make(chan *chess.Game)
	b := make(chan *chess.Game)
	done := make(chan bool, 1)

	go wplay(w, b, done)
	go bplay(w, b, done)
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	g := startGame(fen)
	// The main goroutine starts the ping/pong by sending into the ping channel
	w <- g
	for {
		// Block the main thread until an interrupt
		if <-done {
			break
		}

		//time.Sleep(time.Second)
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

	results, err := engW.GoDepth(20, resultOpts)
	if err != nil {
		fmt.Println(err.Error())
	}
	return results.BestMove
}

func bestMoveBlack(fens string) string {
	engB.SetFEN(fens)

	results, err := engB.GoDepth(20, resultOpts)
	if err != nil {
		fmt.Println(err.Error())
	}
	return results.BestMove
}
