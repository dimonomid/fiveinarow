package main

import (
	"fmt"

	"dmitryfrank.com/fiveinarow/engine"
)

func main() {
	b := engine.MakeBoard(50, 50)
	fmt.Printf("board=%v\n", b)

	for {
		pposx := b.FindBestPos(engine.PointX)
		if pposx == nil {
			fmt.Printf("No more positions!\n")
			return
		}
		b.Put(pposx.X, pposx.Y, engine.PointX)
		fmt.Printf("X: x=%d, y=%d\n", pposx.X, pposx.Y)
		printBoard(b)

		pposo := b.FindBestPos(engine.PointO)
		if pposx == nil {
			fmt.Printf("No more positions!\n")
			return
		}
		b.Put(pposo.X, pposo.Y, engine.PointO)
		fmt.Printf("O: x=%d, y=%d\n", pposo.X, pposo.Y)
		printBoard(b)

	}
}

func printBoard(b *engine.Board) {
	fmt.Println("----------------------")
	for y := range b.Cells {
		for _, c := range b.Cells[y] {
			switch c {
			case engine.CellEmpty:
				fmt.Print(". ")
			case engine.CellX:
				fmt.Print("x ")
			case engine.CellO:
				fmt.Print("o ")
			case engine.CellXStriked:
				fmt.Print("##")
			case engine.CellOStriked:
				fmt.Print("()")
			default:
				panic("hey")
			}
		}
		fmt.Println()
	}
	fmt.Println("----------------------")
}
