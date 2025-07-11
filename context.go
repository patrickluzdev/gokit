package gokit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Ctx struct {
	request  *http.Request
	response http.ResponseWriter
}

func (c *Ctx) Body() []byte {
	body, err := io.ReadAll(c.request.Body)
	if err != nil {
		return nil
	}
	return body
}

func (c *Ctx) Data(code int, data []byte) {
	c.response.WriteHeader(code)
	c.response.Header().Add("Content-Type", "application/json")
	c.response.Write(data)
}

func (c *Ctx) Header(key string) string {
	return c.request.Header.Get(key)
}

func (c *Ctx) JSON(code int, obj any) {
	data, _ := json.Marshal(obj)
	c.response.WriteHeader(code)
	c.response.Header().Add("Content-Type", "application/json")
	c.response.Write(data)
}

func (c *Ctx) Param(key string) string {
	return c.request.PathValue(key)
}

func (c *Ctx) ParseJSON(v any) error {
	body, err := io.ReadAll(c.request.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

func (c *Ctx) ParseString() string {
	body, err := io.ReadAll(c.request.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

func (c *Ctx) Query(key string) string {
	return c.request.URL.Query().Get(key)
}

func (c *Ctx) Request() *http.Request {
	return c.request
}

func (c *Ctx) String(code int, format string, values ...any) {
	body := fmt.Sprintf(format, values...)
	c.response.WriteHeader(code)
	c.response.Header().Add("Content-Type", "text/plain; charset=utf-8")
	c.response.Write([]byte(body))
}

func (c *Ctx) Writer() http.ResponseWriter {
	return c.response
}

func NewCtx(w http.ResponseWriter, r *http.Request) Context {
	return &Ctx{
		request:  r,
		response: w,
	}
}
