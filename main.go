//+build js

package main

import (
	"encoding/json"
	"syscall/js"
)

type PromiseFn func(resolve js.Value, reject js.Value, args ...js.Value) interface{}

func JsError(err error) js.Value {
	return js.Global().Get("Error").New(err.Error())
}

func JSPromise(fn PromiseFn) js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) interface{} {

		handler := js.FuncOf(func(_ js.Value, argsPromise []js.Value) interface{} {
			go func() {
				fn(argsPromise[0], argsPromise[1], args...)
			}()

			return nil
		})

		return js.Global().Get("Promise").New(handler)
	})
}

func PHPForm(_ js.Value, args []js.Value) interface{} {
	form := NewForm()

	if len(args) > 0 && !args[0].IsUndefined() {
		form.SetPrefix(args[0].String())
	}

	return map[string]interface{}{
		"stringify": JSPromise(func(resolve js.Value, reject js.Value, args ...js.Value) interface{} {
			var ms []map[string]interface{}

			err := json.Unmarshal([]byte(args[0].String()), &ms)
			if err != nil {
				return reject.Invoke(JsError(err))
			}

			var fields []Field
			for _, m := range ms {
				fields = append(fields, *NewFieldWithMap(m))
			}

			res, err := form.Stringify(fields)
			if err != nil {
				return reject.Invoke(JsError(err))
			}

			return resolve.Invoke(res)
		}),
		"parse": JSPromise(func(resolve js.Value, reject js.Value, args ...js.Value) interface{} {
			fields, err := form.ParseCode(args[0].String())

			if err != nil {
				return reject.Invoke(JsError(err))
			}

			var res []interface{}
			for _, field := range fields {
				res = append(res, field.ToMap())
			}

			return resolve.Invoke(js.ValueOf(res))
		}),
	}
}

func main() {
	js.Global().Set("PHPFormFunc", js.FuncOf(PHPForm))
	<-make(chan struct{})
}
