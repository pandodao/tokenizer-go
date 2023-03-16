package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pandodao/tokenizer-go"
)

func main() {
	t := flag.String("text", "", "text to calculate token")
	flag.Parse()

	if *t == "" {
		fmt.Println("Please specify text to calculate token")
		os.Exit(1)
	}

	fmt.Println(tokenizer.MustCalToken(*t))
}
