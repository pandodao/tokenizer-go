package tokenizer

import (
	_ "embed"
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

	vm       *goja.Runtime
	calToken goja.Callable
)

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
		"const tokenizer = new GPT3NodeTokenizer({type: 'gpt3'}); function calToken(str) {return tokenizer.encode(str).bpe.length}")
	if err != nil {
		panic(err)
	}

	var ok bool
	calToken, ok = goja.AssertFunction(vm.Get("calToken"))
	if !ok {
		panic("not a function")
	}
}

func MustCalToken(str string) int64 {
	token, err := CalToken(str)
	if err != nil {
		panic(err)
	}
	return token
}

func CalToken(str string) (int64, error) {
	v, err := calToken(goja.Undefined(), vm.ToValue(str))
	if err != nil {
		return 0, err
	}
	return v.ToInteger(), nil
}
