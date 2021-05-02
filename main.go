//+build js

package main

import (
	"encoding/json"
	"syscall/js"
)

type PromiseFn func(resolve js.Value, reject js.Value, args ...js.Value) interface{}

func JsError(err error) js.Error {
	return js.Error{
		Value: js.ValueOf(err.Error()),
	}
}

func JSPromise(fn PromiseFn) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		handler := js.FuncOf(func(this js.Value, argsPromise []js.Value) interface{} {
			go func() {
				fn(argsPromise[0], argsPromise[1], args...)
			}()

			return nil
		})

		return js.Global().Get("Promise").New(handler)
	})
}

func PHPForm(this js.Value, args []js.Value) interface{} {

	form := NewForm(args[0].String())

	if len(args) > 1 {
		form = form.WithPrefix(args[1].String())
	}

	return map[string]interface{}{
		"stringify": JSPromise(func(resolve js.Value, reject js.Value, args ...js.Value) interface{} {
			s := args[0].String()

			var ms []map[string]interface{}

			err := json.Unmarshal([]byte(s), &ms)
			if err != nil {
				return reject.Invoke(JsError(err))
			}

			var inputs []Input
			for _, m := range ms {
				inputs = append(inputs, *NewInputWithMap(m))
			}

			res, err := form.GenerateCodeWithInputs(inputs)
			if err != nil {
				return reject.Invoke(JsError(err))
			}

			return resolve.Invoke(res)
		}),
		"parse": JSPromise(func(resolve js.Value, reject js.Value, args ...js.Value) interface{} {
			inputs, err := form.GenerateInputs()
			if err != nil {
				return reject.Invoke(JsError(err))
			}

			var res []interface{}
			for _, input := range inputs {

				m, err := input.ToMap()
				if err != nil {
					return reject.Invoke(JsError(err))
				}

				res = append(res, m)
			}

			return resolve.Invoke(js.ValueOf(res))
		}),
	}
}

func main() {
	js.Global().Set("PHPForm", js.FuncOf(PHPForm))
	<-make(chan struct{})
}
