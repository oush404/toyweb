package recovery

import "toyweb"

type MiddlewareBuilder struct {
	StatusCode int
	ErrMsg     string
	LogFunc    func(ctx *toyweb.Context)
}

func (m *MiddlewareBuilder) Build() toyweb.Middleware {
	return func(next toyweb.HandleFunc) toyweb.HandleFunc {
		return func(ctx *toyweb.Context) {
			defer func() {
				if err := recover(); err != nil {
					ctx.RespStatusCode = m.StatusCode
					ctx.RespData = []byte(m.ErrMsg)
					m.LogFunc(ctx)
				}
			}()
			next(ctx)
		}
	}
}
