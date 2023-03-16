package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/pandodao/tokenizer-go"
)

func main() {
	token := flag.String("token", "", "text to calculate token")
	encode := flag.String("encode", "", "text to encode")
	decode := flag.String("decode", "", "tokens to decode")
	flag.Parse()

	switch {
	case *token != "":
		fmt.Println(tokenizer.MustCalToken(*token))
	case *encode != "":
		data, _ := json.Marshal(tokenizer.MustEncode(*encode))
		fmt.Println(string(data))
	case *decode != "":
		var s []int
		if err := json.Unmarshal([]byte(*decode), &s); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(tokenizer.MustDecode(s))
	default:
		flag.Usage()
	}
}
