package toyweb

type Middleware func(next HandleFunc) HandleFunc
