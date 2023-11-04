package main

import (
	"C"

	"github.com/valyala/fasthttp"
)
import (
	"unsafe"

	"github.com/tidwall/sjson"
)

func main() {

}

// !STRUCTS

type Generic_Dict struct {
	Key   *C.char
	Value *C.char
}

type Request_c struct {
	Data *C.char
}

type Callback func(Request_c) unsafe.Pointer

var methods_functions = map[string]map[string]Callback{}

//!EXPORTED FUNCTIONS

//export RegisterCallback
func RegisterCallback(id *C.char, fn Callback, method *C.char) {
	funcs := methods_functions[C.GoString(method)]

	if funcs == nil {
		funcs = map[string]Callback{C.GoString(id): fn}
	} else {
		funcs[C.GoString(id)] = fn
	}

	methods_functions[C.GoString(method)] = funcs

}

//export StartServer
func StartServer(host *C.char, port *C.char) {

	m := func(ctx *fasthttp.RequestCtx) {

		funcc, e := methods_functions[string(ctx.Method())][string(ctx.Path())]

		if e {
			data_py, _ := sjson.Set("{}", "URI", string(ctx.Path()))
			data_py, _ = sjson.Set(data_py, "Headers", ctx.Request.Header.String())
			data_py, _ = sjson.Set(data_py, "Body", ctx.Request.Body())

			req_c := Request_c{Data: C.CString(data_py)}
			resp := funcc(req_c)
			ctx.WriteString(C.GoString((*C.char)(resp)))

		} else {
			ctx.WriteString("not found")
		}

	}

	fasthttp.ListenAndServe(C.GoString(host)+":"+C.GoString(port), m)

}
