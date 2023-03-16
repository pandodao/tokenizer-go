package tokenizer

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	tables := []struct {
		input      string
		result     EncodeResult
		ignoreText bool
	}{
		{
			input: "Hello World",
			result: EncodeResult{
				Bpe:  []int{15496, 2159},
				Text: []string{"Hello", " World"},
			},
		},
		{
			input:      "你好，世界",
			ignoreText: true,
			result: EncodeResult{
				Bpe: []int{19526, 254, 25001, 121, 171, 120, 234, 10310, 244, 45911, 234},
			},
		},
	}

	for _, table := range tables {
		t.Run(table.input, func(t *testing.T) {
			start := time.Now()
			r := MustEncode(table.input)
			if !reflect.DeepEqual(r.Bpe, table.result.Bpe) {
				t.Errorf("Encode Bpe was incorrect, got: %v, want: %v.", r.Bpe, table.result.Bpe)
			}
			if !table.ignoreText && !reflect.DeepEqual(r.Text, table.result.Text) {
				t.Errorf("Encode Text was incorrect, got: %v, want: %v.", r.Text, table.result.Text)
			}
			t.Logf("Encode(%s) cost: %s", table.input, time.Since(start))
		})
	}
}

func TestDecode(t *testing.T) {
	tables := []struct {
		input  []int
		result string
	}{
		{
			input:  []int{15496, 2159},
			result: "Hello World",
		},
		{
			input:  []int{19526, 254, 25001, 121, 171, 120, 234, 10310, 244, 45911, 234},
			result: "你好，世界",
		},
	}

	for _, table := range tables {
		t.Run(fmt.Sprintf("%v", table.input), func(t *testing.T) {
			start := time.Now()
			r := MustDecode(table.input)
			if r != table.result {
				t.Errorf("Decode was incorrect, got: %v, want: %v.", r, table.result)
			}
			t.Logf("Decode(%v) cost: %s", table.input, time.Since(start))
		})
	}
}

func TestCalToken(t *testing.T) {
	tables := []struct {
		input string
		token int
	}{
		{
			input: "Hello World",
			token: 2,
		},
		{
			input: "你好，世界",
			token: 11,
		},
	}

	for _, table := range tables {
		t.Run(table.input, func(t *testing.T) {
			start := time.Now()
			token := MustCalToken(table.input)
			if token != table.token {
				t.Errorf("CalToken was incorrect, got: %d, want: %d.", token, table.token)
			}
			t.Logf("CalToken(%s) cost: %s", table.input, time.Since(start))
		})
	}
}

func BenchmarkCalToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MustCalToken(`Many words map to one token, but some don't: indivisible.

Unicode characters like emojis may be split into many tokens containing the underlying bytes: 🤚🏾

Sequences of characters commonly found next to each other may be grouped together: 1234567890`)
	}
}
