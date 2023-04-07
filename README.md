When we attempted to fulfill this requirement, we couldn't find a readily available Go package to use. Due to time constraints, we resorted to a workaround by calling JavaScript instead. However, this approach was not elegant and not very efficient. Now there is a native Go package implementation available at https://github.com/pkoukk/tiktoken-go, please prioritize using it.

---

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
=== RUN   TestNewGojaRuntime
--- PASS: TestNewGojaRuntime (0.00s)
=== RUN   TestValidateFunctionsWithinGojaRuntime
--- PASS: TestValidateFunctionsWithinGojaRuntime (0.61s)
=== RUN   TestEncode
=== RUN   TestEncode/ASCII_Characters
    tokenizer_test.go:117: Encode(Hello World) cost: 620.252292ms
=== RUN   TestEncode/CJK_Characters
    tokenizer_test.go:117: Encode(你好，世界) cost: 387.25µs
=== RUN   TestEncode/WithConcurrency
    tokenizer_test.go:172: Encode(ASCII_Characters) ran 20 times concurrently, cost average: 361.588418ms, cost min: 75.833µs, cost max: 1.829107916s
    tokenizer_test.go:178: Encode(CJK_Characters) ran 20 times concurrently, cost average: 446.462658ms, cost min: 170.292µs, cost max: 1.831984708s
--- PASS: TestEncode (2.45s)
    --- PASS: TestEncode/ASCII_Characters (0.62s)
    --- PASS: TestEncode/CJK_Characters (0.00s)
    --- PASS: TestEncode/WithConcurrency (1.83s)
=== RUN   TestDecode
=== RUN   TestDecode/ASCII_Characters
    tokenizer_test.go:212: Decode([15496 2159]) cost: 150.416µs
=== RUN   TestDecode/CJK_Characters
    tokenizer_test.go:212: Decode([19526 254 25001 121 171 120 234 10310 244 45911 234]) cost: 34.584µs
=== RUN   TestDecode/WithConcurrency
    tokenizer_test.go:258: Decode(ASCII_Characters) ran 20 times concurrently, cost average: 45.558µs, cost min: 29.708µs, cost max: 153.458µs
    tokenizer_test.go:264: Decode(CJK_Characters) ran 20 times concurrently, cost average: 62.145µs, cost min: 37.291µs, cost max: 183.292µs
--- PASS: TestDecode (0.00s)
    --- PASS: TestDecode/ASCII_Characters (0.00s)
    --- PASS: TestDecode/CJK_Characters (0.00s)
    --- PASS: TestDecode/WithConcurrency (0.00s)
=== RUN   TestCalToken
=== RUN   TestCalToken/ASCII_Characters
    tokenizer_test.go:298: CalToken(Hello World) cost: 357.583µs
=== RUN   TestCalToken/CJK_Characters
    tokenizer_test.go:298: CalToken(你好，世界) cost: 217.709µs
=== RUN   TestCalToken/WithConcurrency
    tokenizer_test.go:344: Decode(ASCII_Characters) ran 20 times concurrently, cost average: 32.636206ms, cost min: 96.75µs, cost max: 647.582833ms
    tokenizer_test.go:350: Decode(CJK_Characters) ran 20 times concurrently, cost average: 429.197µs, cost min: 230.375µs, cost max: 1.167416ms
--- PASS: TestCalToken (0.65s)
    --- PASS: TestCalToken/ASCII_Characters (0.00s)
    --- PASS: TestCalToken/CJK_Characters (0.00s)
    --- PASS: TestCalToken/WithConcurrency (0.65s)
goos: darwin
goarch: arm64
pkg: github.com/pandodao/tokenizer-go
BenchmarkCalToken
BenchmarkCalToken/ASCII_Characters
BenchmarkCalToken/ASCII_Characters-10                546           2186558 ns/op
BenchmarkCalToken/CJK_Characters
BenchmarkCalToken/CJK_Characters-10                  420           2942631 ns/op
PASS
ok      github.com/pandodao/tokenizer-go        10.869s
```

## Thanks

* https://github.com/botisan-ai/gpt3-tokenizer
* https://github.com/dop251/goja

## License
See the [LICENSE](https://github.com/pandodao/tokenizer-go/blob/main/LICENSE) file.
