package tokenizer

import (
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func max[T int | time.Duration](slice []T) T {
	sort.SliceStable(slice, func(i, j int) bool {
		return slice[i] < slice[j]
	})

	return slice[len(slice)-1]
}

func min[T int | time.Duration](slice []T) T {
	sort.SliceStable(slice, func(i, j int) bool {
		return slice[i] < slice[j]
	})

	return slice[0]
}

func average[T int | time.Duration](slice []T) T {
	var sum T
	for _, item := range slice {
		sum += item
	}

	return sum / T(len(slice))
}

func TestNewGojaRuntime(t *testing.T) {
	originalTokenizerJs := tokenizerJs
	defer func() {
		tokenizerJs = originalTokenizerJs
	}()

	tokenizerJs = ""
	runtime := newGojaRuntime()
	require.Error(t, runtime.err)
	assert.EqualError(t, runtime.err, "ReferenceError: GPT3NodeTokenizer is not defined at <eval>:2:23(3)")
}

func TestValidateFunctionsWithinGojaRuntime(t *testing.T) {
	vm := goja.New()
	registry.Enable(vm)

	encode, decode, err := getEncodeAndDecodeFunctionsWithinGojaRuntime(vm)
	require.Error(t, err)
	assert.EqualError(t, err, "encode is not a function")
	assert.Nil(t, encode)
	assert.Nil(t, decode)

	_, err = vm.RunString(tokenizerJs + "\n" +
		`const tokenizer = new GPT3NodeTokenizer({type: 'gpt3'});
		 function encode(str) {return tokenizer.encode(str)}`)
	require.NoError(t, err)

	encode, decode, err = getEncodeAndDecodeFunctionsWithinGojaRuntime(vm)
	require.Error(t, err)
	assert.EqualError(t, err, "decode is not a function")
	assert.Nil(t, encode)
	assert.Nil(t, decode)

	_, err = vm.RunString("function decode(tokens) {return tokenizer.decode(tokens)}")
	require.NoError(t, err)

	encode, decode, err = getEncodeAndDecodeFunctionsWithinGojaRuntime(vm)
	require.NoError(t, err)
	assert.NotNil(t, encode)
	assert.NotNil(t, decode)
}

type testEncode struct {
	testName   string
	input      string
	result     EncodeResult
	ignoreText bool
}

func TestEncode(t *testing.T) {
	tables := []testEncode{
		{
			testName: "ASCII_Characters",
			input:    "Hello World",
			result: EncodeResult{
				Bpe:  []int{15496, 2159},
				Text: []string{"Hello", " World"},
			},
		},
		{
			testName:   "CJK_Characters",
			input:      "你好，世界",
			ignoreText: true,
			result: EncodeResult{
				Bpe: []int{19526, 254, 25001, 121, 171, 120, 234, 10310, 244, 45911, 234},
			},
		},
	}

	var ignoreTextRan bool
	for _, table := range tables {
		t.Run(table.testName, func(t *testing.T) {
			start := time.Now()
			r := MustEncode(table.input)
			assert.Equal(t, table.result.Bpe, r.Bpe)
			if !table.ignoreText {
				assert.Equal(t, table.result.Text, r.Text)
				ignoreTextRan = true
			}

			t.Logf("Encode(%s) cost: %s", table.input, time.Since(start))
		})
	}
	assert.True(t, ignoreTextRan)
	ignoreTextRan = false

	t.Run("WithConcurrency", func(t *testing.T) {
		concurrency := 20

		tablesMat := make([][]testEncode, concurrency)
		resultsMat := make([][]EncodeResult, concurrency)
		timeCostsMat := make([][]time.Duration, concurrency)
		for i := range tablesMat {
			tablesMat[i] = tables
			resultsMat[i] = make([]EncodeResult, len(tables))    // init
			timeCostsMat[i] = make([]time.Duration, len(tables)) // init
		}

		var wg sync.WaitGroup
		for i, elem := range tablesMat {
			for j := range elem {
				wg.Add(1)
				go func(iIndex, jIndex int) {
					start := time.Now()
					table := tablesMat[iIndex][jIndex]
					result := MustEncode(table.input)

					resultsMat[iIndex][jIndex] = result
					timeCostsMat[iIndex][jIndex] = time.Since(start)
					wg.Done()
				}(i, j)
			}
		}
		wg.Wait()

		for i, ts := range tablesMat {
			for j := range ts {
				r := resultsMat[i][j]
				assert.Equal(t, ts[j].result.Bpe, r.Bpe)
				if !ts[j].ignoreText {
					assert.Equal(t, ts[j].result.Text, r.Text)
					ignoreTextRan = true
				}
			}
		}

		assert.True(t, ignoreTextRan)

		timeCostsForASCIICharacters := make([]time.Duration, len(timeCostsMat))
		timeCostsForCJKCharacters := make([]time.Duration, len(timeCostsMat))
		for i := range timeCostsMat {
			timeCostsForASCIICharacters[i] = timeCostsMat[i][0]
			timeCostsForCJKCharacters[i] = timeCostsMat[i][1]
		}

		t.Logf("Encode(ASCII_Characters) ran %d times concurrently, cost average: %s, cost min: %s, cost max: %s",
			concurrency,
			average(timeCostsForASCIICharacters),
			min(timeCostsForASCIICharacters),
			max(timeCostsForASCIICharacters),
		)
		t.Logf("Encode(CJK_Characters) ran %d times concurrently, cost average: %s, cost min: %s, cost max: %s",
			concurrency,
			average(timeCostsForCJKCharacters),
			min(timeCostsForCJKCharacters),
			max(timeCostsForCJKCharacters),
		)
	})
}

type testDecode struct {
	testName string
	input    []int
	result   string
}

func TestDecode(t *testing.T) {
	tables := []testDecode{
		{
			testName: "ASCII_Characters",
			input:    []int{15496, 2159},
			result:   "Hello World",
		},
		{
			testName: "CJK_Characters",
			input:    []int{19526, 254, 25001, 121, 171, 120, 234, 10310, 244, 45911, 234},
			result:   "你好，世界",
		},
	}

	for _, table := range tables {
		t.Run(table.testName, func(t *testing.T) {
			start := time.Now()
			r := MustDecode(table.input)
			assert.Equal(t, table.result, r)
			t.Logf("Decode(%v) cost: %s", table.input, time.Since(start))
		})
	}

	t.Run("WithConcurrency", func(t *testing.T) {
		concurrency := 20

		tablesMat := make([][]testDecode, concurrency)
		resultsMat := make([][]string, concurrency)
		timeCostsMat := make([][]time.Duration, concurrency)
		for i := range tablesMat {
			tablesMat[i] = tables
			resultsMat[i] = make([]string, len(tables))          // init
			timeCostsMat[i] = make([]time.Duration, len(tables)) // init
		}

		var wg sync.WaitGroup
		for i, elem := range tablesMat {
			for j := range elem {
				wg.Add(1)
				go func(iIndex, jIndex int) {
					start := time.Now()
					table := tablesMat[iIndex][jIndex]
					result := MustDecode(table.input)

					resultsMat[iIndex][jIndex] = result
					timeCostsMat[iIndex][jIndex] = time.Since(start)
					wg.Done()
				}(i, j)
			}
		}
		wg.Wait()

		for i, elem := range tablesMat {
			for j := range elem {
				r := resultsMat[i][j]
				assert.Equal(t, elem[j].result, r)
			}
		}

		timeCostsForASCIICharacters := make([]time.Duration, len(timeCostsMat))
		timeCostsForCJKCharacters := make([]time.Duration, len(timeCostsMat))
		for i := range timeCostsMat {
			timeCostsForASCIICharacters[i] = timeCostsMat[i][0]
			timeCostsForCJKCharacters[i] = timeCostsMat[i][1]
		}

		t.Logf("Decode(ASCII_Characters) ran %d times concurrently, cost average: %s, cost min: %s, cost max: %s",
			concurrency,
			average(timeCostsForASCIICharacters),
			min(timeCostsForASCIICharacters),
			max(timeCostsForASCIICharacters),
		)
		t.Logf("Decode(CJK_Characters) ran %d times concurrently, cost average: %s, cost min: %s, cost max: %s",
			concurrency,
			average(timeCostsForCJKCharacters),
			min(timeCostsForCJKCharacters),
			max(timeCostsForCJKCharacters),
		)
	})
}

type testCalToken struct {
	testName string
	input    string
	token    int
}

func TestCalToken(t *testing.T) {
	tables := []testCalToken{
		{
			testName: "ASCII_Characters",
			input:    "Hello World",
			token:    2,
		},
		{
			testName: "CJK_Characters",
			input:    "你好，世界",
			token:    11,
		},
	}

	for _, table := range tables {
		t.Run(table.testName, func(t *testing.T) {
			start := time.Now()
			token := MustCalToken(table.input)
			assert.Equal(t, table.token, token)
			t.Logf("CalToken(%s) cost: %s", table.input, time.Since(start))
		})
	}

	t.Run("WithConcurrency", func(t *testing.T) {
		concurrency := 20

		tablesMat := make([][]testCalToken, concurrency)
		resultsMat := make([][]int, concurrency)
		timeCostsMat := make([][]time.Duration, concurrency)
		for i := range tablesMat {
			tablesMat[i] = tables
			resultsMat[i] = make([]int, len(tables))             // init
			timeCostsMat[i] = make([]time.Duration, len(tables)) // init
		}

		var wg sync.WaitGroup
		for i, elem := range tablesMat {
			for j := range elem {
				wg.Add(1)
				go func(iIndex, jIndex int) {
					start := time.Now()
					table := tablesMat[iIndex][jIndex]
					result := MustCalToken(table.input)

					resultsMat[iIndex][jIndex] = result
					timeCostsMat[iIndex][jIndex] = time.Since(start)
					wg.Done()
				}(i, j)
			}
		}
		wg.Wait()

		for i, elem := range tablesMat {
			for j := range elem {
				token := resultsMat[i][j]
				assert.Equal(t, elem[j].token, token)
			}
		}

		timeCostsForASCIICharacters := make([]time.Duration, len(timeCostsMat))
		timeCostsForCJKCharacters := make([]time.Duration, len(timeCostsMat))
		for i := range timeCostsMat {
			timeCostsForASCIICharacters[i] = timeCostsMat[i][0]
			timeCostsForCJKCharacters[i] = timeCostsMat[i][1]
		}

		t.Logf("Decode(ASCII_Characters) ran %d times concurrently, cost average: %s, cost min: %s, cost max: %s",
			concurrency,
			average(timeCostsForASCIICharacters),
			min(timeCostsForASCIICharacters),
			max(timeCostsForASCIICharacters),
		)
		t.Logf("Decode(CJK_Characters) ran %d times concurrently, cost average: %s, cost min: %s, cost max: %s",
			concurrency,
			average(timeCostsForCJKCharacters),
			min(timeCostsForCJKCharacters),
			max(timeCostsForCJKCharacters),
		)
	})
}

func BenchmarkCalToken(b *testing.B) {
	b.Run("ASCII_Characters", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = MustCalToken(`Many words map to one token, but some don't: indivisible.

Unicode characters like emojis may be split into many tokens containing the underlying bytes: 🤚🏾

Sequences of characters commonly found next to each other may be grouped together: 1234567890`)
		}
	})

	b.Run("CJK_Characters", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = MustCalToken(`许多词都会被映射到一个令牌上，但有些词的类型不会：不可分割的。

像 Emoji 这样的 Unicode 字符可以被分割成许多包含底层字节的标记：🤚🏾

常见的字符序列彼此相邻，可以归为一组：1234567890`)
		}
	})
}
