package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"strings"
	"unsafe"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/valyala/fasthttp"
)

func main() {

}

// !STRUCTS

type Generic_Dict struct {
	Key   *C.char
	Value *C.char
}

type Request_c struct {
	Data         *C.char
	Body_Address *C.uchar
	Body_Len     C.int
}

type Callback func(*Request_c) unsafe.Pointer

var methods_functions = map[string]map[string]Callback{}
var static_paths = map[string]fasthttp.RequestHandler{}
var error_page = "Internal Server Error"
var not_found_page = "Not Found"

//!EXPORTED FUNCTIONS

//export MountStatic
func MountStatic(path *C.char, directory *C.char) {
	fs := &fasthttp.FS{
		Root: C.GoString(directory),
	}
	fs.PathRewrite = fasthttp.NewPathSlashesStripper(1)
	static_paths[C.GoString(path)] = fs.NewRequestHandler()

}

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

//export SetNotFoundPage
func SetNotFoundPage(html *C.char) {
	not_found_page = C.GoString(html)
}

//export SetErrorPage
func SetErrorPage(html *C.char) {
	error_page = C.GoString(html)
}

//export StartServer
func StartServer(host *C.char, port *C.char) {

	m := func(ctx *fasthttp.RequestCtx) {

		for path, static := range static_paths {
			if strings.HasPrefix(string(ctx.Path()), path) {
				static(ctx)
				return
			}
		}

		funcc, e := methods_functions[string(ctx.Method())][string(ctx.Path())]

		if e {
			data_py, _ := sjson.Set("{}", "URI", string(ctx.Path()))
			data_py, _ = sjson.Set(data_py, "Headers", ctx.Request.Header.String())
			data_py, _ = sjson.Set(data_py, "Method", string(ctx.Method()))

			data_c := C.CString(data_py)
			req_c := Request_c{Data: data_c, Body_Len: C.int(0)}

			body := ctx.Request.Body()
			body_len := len(body)

			if body_len != 0 {
				req_c.Body_Address = (*C.uchar)(unsafe.Pointer(&body[0]))
				req_c.Body_Len = C.int(body_len)
			}

			resp := funcc(&req_c)
			C.free(unsafe.Pointer(data_c))

			resp_go := gjson.Parse(C.GoString((*C.char)(resp)))

			if !resp_go.Get("status_code").Exists() {
				ctx.WriteString(error_page)
				ctx.SetStatusCode(500)
				return
			}

			ctx.WriteString(resp_go.Get("content").Str)
			ctx.SetStatusCode(int(resp_go.Get("status_code").Int()))

			for k, v := range resp_go.Get("headers").Map() {
				ctx.Response.Header.Add(k, v.String())
			}

			for _, v := range resp_go.Get("cookies").Array() {
				coo := fasthttp.Cookie{}

				coo.SetKey(v.Get("key").String())
				coo.SetValue(v.Get("value").String())

				if v.Get("expire").Exists() {
					coo.SetExpire(v.Get("expire").Time())
				}
				if v.Get("domain").Exists() {
					coo.SetDomain(v.Get("domain").String())
				}
				if v.Get("maxAge").Exists() {
					coo.SetMaxAge(int(v.Get("maxAge").Int()))
				}

				ctx.Response.Header.Cookie(&coo)
			}

		} else {
			ctx.WriteString(not_found_page)
			ctx.SetStatusCode(404)
			return
		}

	}

	fasthttp.ListenAndServe(C.GoString(host)+":"+C.GoString(port), m)

}
