package errhdl

import (
	"bytes"
	"html/template"
	"testing"
	"toyweb"
)

func TestNewMiddlewareBuilder(t *testing.T) {
	s := toyweb.NewHTTPServer()
	s.Get("/user", func(ctx *toyweb.Context) {
		ctx.RespData = []byte("hello,world")
	})

	page := `
<html>
	<h1>404 NOT FOUND</h1>
</html>
`

	tpl, err := template.New("404").Parse(page)
	if err != nil {
		t.Fatal(err)
	}
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, nil)
	if err != nil {
		t.Fatal(err)
	}
	s.Use(NewMiddlewareBuilder().RegisterError(404, buffer.Bytes()).Build())
	s.Start(":8080")
}
