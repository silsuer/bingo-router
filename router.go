package bingo_router

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"fmt"
)

type Router struct {
	*httprouter.Router
}

func New() *Router {
	return &Router{
		httprouter.New(),
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	fmt.Fprint(w,"aaa")

	// 全局中间件
	//r.Router.ServeHTTP(w, req)
	//if r.PanicHandler != nil {
	//	defer r.recv(w, req)
	//}
	//
	//path := req.URL.Path
	//
	//if root := r.trees[req.Method]; root != nil {
	//	if handle, ps, tsr := root.getValue(path); handle != nil {
	//		handle(w, req, ps)
	//		return
	//	} else if req.Method != "CONNECT" && path != "/" {
	//		code := 301 // Permanent redirect, request with GET method
	//		if req.Method != "GET" {
	//			// Temporary redirect, request with same method
	//			// As of Go 1.3, Go does not support status code 308.
	//			code = 307
	//		}
	//
	//		if tsr && r.RedirectTrailingSlash {
	//			if len(path) > 1 && path[len(path)-1] == '/' {
	//				req.URL.Path = path[:len(path)-1]
	//			} else {
	//				req.URL.Path = path + "/"
	//			}
	//			http.Redirect(w, req, req.URL.String(), code)
	//			return
	//		}
	//
	//		// Try to fix the request path
	//		if r.RedirectFixedPath {
	//			fixedPath, found := root.findCaseInsensitivePath(
	//				CleanPath(path),
	//				r.RedirectTrailingSlash,
	//			)
	//			if found {
	//				req.URL.Path = string(fixedPath)
	//				http.Redirect(w, req, req.URL.String(), code)
	//				return
	//			}
	//		}
	//	}
	//}
	//
	//if req.Method == "OPTIONS" && r.HandleOPTIONS {
	//	// Handle OPTIONS requests
	//	if allow := r.allowed(path, req.Method); len(allow) > 0 {
	//		w.Header().Set("Allow", allow)
	//		return
	//	}
	//} else {
	//	// Handle 405
	//	if r.HandleMethodNotAllowed {
	//		if allow := r.allowed(path, req.Method); len(allow) > 0 {
	//			w.Header().Set("Allow", allow)
	//			if r.MethodNotAllowed != nil {
	//				r.MethodNotAllowed.ServeHTTP(w, req)
	//			} else {
	//				http.Error(w,
	//					http.StatusText(http.StatusMethodNotAllowed),
	//					http.StatusMethodNotAllowed,
	//				)
	//			}
	//			return
	//		}
	//	}
	//}
	//
	//// Handle 404
	//if r.NotFound != nil {
	//	r.NotFound.ServeHTTP(w, req)
	//} else {
	//	http.NotFound(w, req)
	//}
}
