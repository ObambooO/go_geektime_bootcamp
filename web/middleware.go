package web

// Middleware 函数式的责任链模式，函数式的洋葱模式
type Middleware func(next HandleFunc) HandleFunc

//type MiddlewareV1 interface {
//	Invoke(next HandleFunc) HandleFunc
//}
//
//type Interceptor interface {
//	Before(ctx *Context)
//	After(ctx *Context)
//	Surround(ctx *Context)
//}

//type Chain []HandleFunc

//type HandleFuncV1 func(ctx *Context) (next bool)
//type ChainV1 struct {
//	handlers []HandleFuncV1
//}
//
//func (c ChainV1) Run(ctx *Context) {
//	for _, h := range c.handlers {
//		next := h(ctx)
//		if !next {
//			return
//		}
//	}
//}
