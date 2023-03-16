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
	fmt.Println(t)
}
```

* As a command line program
```
~ % tokenizer -text "Many words map to one token, but some don't: indivisible.

Unicode characters like emojis may be split into many tokens containing the underlying bytes: 🤚🏾

Sequences of characters commonly found next to each other may be grouped together: 1234567890"
64
~ %
```

## Benchmark

```
% go test -v -bench=.
=== RUN   TestCalToken
    tokenizer_test.go:29: CalToken(Hello World) cost: 954.578µs
    tokenizer_test.go:29: CalToken(你好，世界) cost: 994.442µs
--- PASS: TestCalToken (0.00s)
goos: darwin
goarch: amd64
pkg: github.com/pandodao/tokenizer-go
cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
BenchmarkCalToken
BenchmarkCalToken-12                 330           3853708 ns/op
PASS
ok      github.com/pandodao/tokenizer-go        2.842s
```

## Thanks

* https://github.com/botisan-ai/gpt3-tokenizer
* https://github.com/dop251/goja

## License
See the [LICENSE](https://github.com/pandodao/tokenizer-go/blob/main/LICENSE) file.
