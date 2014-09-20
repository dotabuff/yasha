package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dotabuff/yasha"
	"github.com/dotabuff/yasha/dota"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected a .dem file as argument")
	}

	for _, path := range os.Args[1:] {
		parser := yasha.ParserFromFile(path)
		parser.OnSayText2 = func(tick int, obj *dota.CUserMsg_SayText2) {
			fmt.Printf("%s - %07d | %s: %s\n", filepath.Base(path), tick, obj.GetPrefix(), obj.GetText())
		}
		parser.Parse()
	}
}
