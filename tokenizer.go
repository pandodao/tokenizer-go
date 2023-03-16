package tokenizer

import (
	_ "embed"
	"encoding/json"
	"path"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

var (
	//go:embed js/gpt3-tokenizer.cjs.development.js
	tokenizerJs string

	//go:embed js/array-keyed-map.js
	arrayKeyedMapJs string

	//go:embed js/text.min.js
	fastTextEncodingJs string

	vm     *goja.Runtime
	encode goja.Callable
	decode goja.Callable
)

type EncodeResult struct {
	Bpe  []int    `json:"bpe"`
	Text []string `json:"text"`
}

func init() {
	vm = goja.New()
	registry := require.NewRegistry(require.WithLoader(func(p string) ([]byte, error) {
		switch path.Base(p) {
		case "array-keyed-map":
			return []byte(arrayKeyedMapJs), nil
		case "fast-text-encoding":
			return []byte(fastTextEncodingJs), nil
		}
		return nil, require.IllegalModuleNameError
	}))

	registry.Enable(vm)
	_, err := vm.RunString(tokenizerJs + "\n" +
		`const tokenizer = new GPT3NodeTokenizer({type: 'gpt3'});
		 function encode(str) {return tokenizer.encode(str)}
		 function decode(tokens) {return tokenizer.decode(tokens)}`)
	if err != nil {
		panic(err)
	}

	var ok bool
	encode, ok = goja.AssertFunction(vm.Get("encode"))
	if !ok {
		panic("encode is not a function")
	}
	decode, ok = goja.AssertFunction(vm.Get("decode"))
	if !ok {
		panic("decode is not a function")
	}
}

func MustCalToken(str string) int {
	token, err := CalToken(str)
	if err != nil {
		panic(err)
	}
	return token
}

func CalToken(str string) (int, error) {
	r, err := Encode(str)
	if err != nil {
		return 0, err
	}

	return len(r.Bpe), nil
}

func MustEncode(str string) EncodeResult {
	r, err := Encode(str)
	if err != nil {
		panic(err)
	}

	return *r
}

func Encode(str string) (*EncodeResult, error) {
	v, err := encode(goja.Undefined(), vm.ToValue(str))
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(v.Export())
	r := &EncodeResult{}
	if err := json.Unmarshal(data, r); err != nil {
		return nil, err
	}

	return r, nil
}

func MustDecode(tokens []int) string {
	r, err := Decode(tokens)
	if err != nil {
		panic(err)
	}

	return r
}

func Decode(tokens []int) (string, error) {
	v, err := decode(goja.Undefined(), vm.ToValue(tokens))
	if err != nil {
		return "", err
	}

	return v.String(), nil
}
