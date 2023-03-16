package tokenizer

import (
	"testing"
	"time"
)

func TestCalToken(t *testing.T) {
	tables := []struct {
		input string
		token int64
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
		start := time.Now()
		token := MustCalToken(table.input)
		if token != table.token {
			t.Errorf("CalToken was incorrect, got: %d, want: %d.", token, table.token)
		}
		t.Logf("CalToken(%s) cost: %s", table.input, time.Since(start))
	}
}

func BenchmarkCalToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MustCalToken(`Many words map to one token, but some don't: indivisible.

Unicode characters like emojis may be split into many tokens containing the underlying bytes: 🤚🏾

Sequences of characters commonly found next to each other may be grouped together: 1234567890`)
	}
}
