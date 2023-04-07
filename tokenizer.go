package tokenizer

import (
	_ "embed"
	"encoding/json"
	"errors"
	"path"
	"sync"

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

	registry *require.Registry

	// optimize the alloc and instancing performance of
	// *goja.Runtime
	pool sync.Pool = sync.Pool{
		New: func() any {
			return newGojaRuntime()
		},
	}
)

// gojaRuntime is a wrapper of *goja.Runtime with the error
type gojaRuntime struct {
	// runtime itself
	vm *goja.Runtime
	// encode function that registered in the runtime
	encode goja.Callable
	// decode function that registered in the runtime
	decode goja.Callable

	// err is the error occurred during the initialization
	err error
}

type EncodeResult struct {
	Bpe  []int    `json:"bpe"`
	Text []string `json:"text"`
}

func init() {
	registry = require.NewRegistry(require.WithLoader(func(p string) ([]byte, error) {
		switch path.Base(p) {
		case "array-keyed-map":
			return []byte(arrayKeyedMapJs), nil
		case "fast-text-encoding":
			return []byte(fastTextEncodingJs), nil
		}
		return nil, require.IllegalModuleNameError
	}))

	// pre-alloc the *goja.Runtime once
	runtime := pool.Get().(*gojaRuntime)
	if runtime.err != nil {
		panic(runtime.err)
	}

	pool.Put(runtime) // put it back to the pool
}

// newGojaRuntime create a new *goja.Runtime and declare the
// tokenizer functions, it returns the wrapped *gojaRuntime with
// the error if any occurred during the initialization
func newGojaRuntime() *gojaRuntime {
	vm := goja.New()
	registry.Enable(vm)
	_, err := vm.RunString(tokenizerJs + "\n" +
		`const tokenizer = new GPT3NodeTokenizer({type: 'gpt3'});
		 function encode(str) {return tokenizer.encode(str)}
		 function decode(tokens) {return tokenizer.decode(tokens)}`)
	if err != nil {
		return &gojaRuntime{
			vm:  vm,
			err: err,
		}
	}

	encode, decode, err := getEncodeAndDecodeFunctionsWithinGojaRuntime(vm)
	return &gojaRuntime{
		vm:     vm,
		encode: encode,
		decode: decode,
		err:    err,
	}
}

// getEncodeAndDecodeFunctionsWithinGojaRuntime returns the encode and
// decode functions within the *goja.Runtime
func getEncodeAndDecodeFunctionsWithinGojaRuntime(vm *goja.Runtime) (goja.Callable, goja.Callable, error) {
	encode, ok := goja.AssertFunction(vm.Get("encode"))
	if !ok {
		return nil, nil, errors.New("encode is not a function")
	}
	decode, ok := goja.AssertFunction(vm.Get("decode"))
	if !ok {
		return nil, nil, errors.New("decode is not a function")
	}

	return encode, decode, nil
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
	gojaRuntime := pool.Get().(*gojaRuntime)
	if gojaRuntime.err != nil {
		return nil, gojaRuntime.err
	}
	defer pool.Put(gojaRuntime) // put it back to the pool

	v, err := gojaRuntime.encode(goja.Undefined(), gojaRuntime.vm.ToValue(str))
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
	gojaRuntime := pool.Get().(*gojaRuntime)
	if gojaRuntime.err != nil {
		return "", gojaRuntime.err
	}
	defer pool.Put(gojaRuntime) // put it back to the pool

	v, err := gojaRuntime.decode(goja.Undefined(), gojaRuntime.vm.ToValue(tokens))
	if err != nil {
		return "", err
	}

	return v.String(), nil
}
