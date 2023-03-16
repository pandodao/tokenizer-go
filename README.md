# tokenizer-go

tokenizer-go is a Go package that simplifies token calculation for OpenAI API users. Although OpenAI does not provide a native Go package for token calculation, tokenizer-go fills the gap by embedding an implementation of an npm package and extracting the results through JavaScript calls. This allows you to use tokenizer-go just like any other Go package in your projects, making it easier to work with token calculations in the Go programming language.

## Install

```shell
# Use as a module
go get -u github.com/pandodao/tokenizer-go

# Use as a command line program
go install  github.com/pandodao/tokenizer-go/cmd/tokenizer@latest
```

## Usage

* As a module
```go
package main

import (
	"fmt"

	"github.com/pandodao/tokenizer-go"
)

func main() {
	t := tokenizer.MustCalToken(`Many words map to one token, but some don't: indivisible.

Unicode characters like emojis may be split into many tokens containing the underlying bytes: 🤚🏾

Sequences of characters commonly found next to each other may be grouped together: 1234567890`)
	fmt.Println(t) // Output: 64

	// Output: {Bpe:[7085 2456 3975 284 530 11241] Text:[Many  words  map  to  one  token]}
	fmt.Printf("%+v\n", tokenizer.MustEncode("Many words map to one token"))

	// Output: Many words map to one token
	fmt.Println(tokenizer.MustDecode([]int{7085, 2456, 3975, 284, 530, 11241}))
}
```

* As a command line program
```
~ % tokenizer -token "hello world"
2
~ %
~ % tokenizer -encode "hello world"
{"bpe":[31373,995],"text":["hello"," world"]}
~ %
~ % tokenizer -decode "[31373,995]"
hello world
~ %
~ % tokenizer
Usage of tokenizer:
  -decode string
        tokens to decode
  -encode string
        text to encode
  -token string
        text to calculate token
~ %
```

## Benchmark

```
% go test -v -bench=.
=== RUN   TestEncode
=== RUN   TestEncode/Hello_World
    tokenizer_test.go:42: Encode(Hello World) cost: 1.151195ms
=== RUN   TestEncode/你好，世界
    tokenizer_test.go:42: Encode(你好，世界) cost: 1.003894ms
--- PASS: TestEncode (0.00s)
    --- PASS: TestEncode/Hello_World (0.00s)
    --- PASS: TestEncode/你好，世界 (0.00s)
=== RUN   TestDecode
=== RUN   TestDecode/[15496_2159]
    tokenizer_test.go:69: Decode([15496 2159]) cost: 124.855µs
=== RUN   TestDecode/[19526_254_25001_121_171_120_234_10310_244_45911_234]
    tokenizer_test.go:69: Decode([19526 254 25001 121 171 120 234 10310 244 45911 234]) cost: 251.501µs
--- PASS: TestDecode (0.00s)
    --- PASS: TestDecode/[15496_2159] (0.00s)
    --- PASS: TestDecode/[19526_254_25001_121_171_120_234_10310_244_45911_234] (0.00s)
=== RUN   TestCalToken
=== RUN   TestCalToken/Hello_World
    tokenizer_test.go:96: CalToken(Hello World) cost: 293.461µs
=== RUN   TestCalToken/你好，世界
    tokenizer_test.go:96: CalToken(你好，世界) cost: 584.905µs
--- PASS: TestCalToken (0.00s)
    --- PASS: TestCalToken/Hello_World (0.00s)
    --- PASS: TestCalToken/你好，世界 (0.00s)
goos: darwin
goarch: amd64
pkg: github.com/pandodao/tokenizer-go
BenchmarkCalToken
BenchmarkCalToken-12                 319           3595615 ns/op
PASS
ok      github.com/pandodao/tokenizer-go        2.833s
```

## Thanks

* https://github.com/botisan-ai/gpt3-tokenizer
* https://github.com/dop251/goja

## License
See the [LICENSE](https://github.com/pandodao/tokenizer-go/blob/main/LICENSE) file.
