package toyweb

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string

	RespStatusCode int
	RespData       []byte

	// 命中的路由
	MatchedRoute string

	cacheQueryValues url.Values

	UserValues map[string]any
}

func (c *Context) BindJSON(val any) error {
	if c.Req.Body == nil {
		return errors.New("web: body为nil")
	}
	decoder := json.NewDecoder(c.Req.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(val)
}

type StringValue struct {
	val string
	err error
}

func (c *Context) FormValue(key string) StringValue {
	if err := c.Req.ParseForm(); err != nil {
		return StringValue{err: err}
	}
	return StringValue{val: c.Req.FormValue(key)}
}

func (c *Context) QueryValue(key string) StringValue {
	if c.cacheQueryValues == nil {
		c.cacheQueryValues = c.Req.URL.Query()
	}
	vals, ok := c.cacheQueryValues[key]
	if !ok {
		return StringValue{err: errors.New("Web: 找不到这个key")}
	}
	return StringValue{val: vals[0]}
}

func (c *Context) PathValue(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{err: errors.New("web:找不到这key")}
	}
	return StringValue{val: val}
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Resp, cookie)
}

func (c *Context) RespJSONOK(val any) error {
	return c.RespJSON(http.StatusOK, val)
}

func (c *Context) RespJSON(code int, val any) error {
	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.Resp.WriteHeader(code)
	_, err = c.Resp.Write(bs)
	return err
}

func (s StringValue) String() (string, error) {
	return s.val, s.err
}

func (s StringValue) ToInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}
