package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Req         *http.Request
	Resp        http.ResponseWriter
	PathParams  map[string]string
	queryValues url.Values

	MatchedRoute string
	//cookieSamSite http.SameSite
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	// 不推荐
	//cookie.SameSite = c.cookieSamSite
	http.SetCookie(c.Resp, cookie)
}

func (c *Context) RespJSONOK(val any) error {
	return c.RespJSON(http.StatusOK, val)
}

func (c *Context) RespJSON(status int, val any) error {
	data, err := json.Marshal(val)
	if err == nil {
		return err
	}
	c.Resp.WriteHeader(status)
	//c.Resp.Header().Set("Content-Type", "application/json")
	//c.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))
	n, err := c.Resp.Write(data)
	if n != len(data) {
		return errors.New("web: 未写入全部数据")
	}
	return err
}

func (c *Context) BindJSON(val any) error {
	if val == nil {
		return errors.New("web: 输入不能为nil")
	}
	if c.Req.Body == nil {
		return errors.New("web: body 不能为nil")
	}
	// 不要用这种写法
	//bs, _ := io.ReadAll(c.Req.Body)
	//json.Unmarshal(bs, val)

	decoder := json.NewDecoder(c.Req.Body)

	// 下面两种不需要支持，在应用层面支持
	// useNumber => 数字会用Number表示
	// 否则会用float64表示
	//decoder.UseNumber()
	// 如果要是有一个未知的字段，就会报错
	// 比如user只有name和id两个字段
	// JSON里面多了一个Age字段，就会报错
	//decoder.DisallowUnknownFields()

	return decoder.Decode(val)
}

// Form 包含query这些其他的参数，所有的表单数据都能拿到
// PostForm 则只包含PATCH、POST、PUT这些请求的参数，只有在编码为x-www-form-urlencoded的时候才能拿到

// FormValue(key1)、FormValue(key2)只会调用一次
func (c *Context) FromValue(key string) StringValue {
	err := c.Req.ParseForm()
	if err != nil {
		return StringValue{
			val: "",
			err: errors.New("web: parse form failed"),
		}
	}
	return StringValue{
		val: c.Req.FormValue(key),
		err: nil,
	}
	// vals是切片
	//vals, ok := c.Req.Form[key]
	//if !ok {
	//	return "", errors.New("web: key not found")
	//}
	//return vals[0], nil
}

// Query 和 表单比起来，它没有缓存
func (c *Context) QueryValue(key string) StringValue {
	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}

	vals, ok := c.queryValues[key]
	if !ok {
		return StringValue{
			val: "",
			err: errors.New("web: key not found"),
		}
	}
	return StringValue{
		val: vals[0],
		err: nil,
	}
	// 用户区别不出来是真的有值，但是值恰好是空字符串还是没有值
	//return c.queryValues.Get(key), nil
}

func (c *Context) PathValue1(key string) (string, error) {
	val, ok := c.PathParams[key]
	if !ok {
		return "", errors.New("web: key not found")
	}
	return val, nil
}

// 达到链式调用的效果
type StringValue struct {
	val string
	err error
}

// 为了能够支持多种返回格式，可引入结构体进行处理
func (c *Context) PathValue(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{
			val: "",
			err: errors.New("web: key not found"),
		}
	}
	return StringValue{
		val: val,
		err: nil,
	}
}

func (s StringValue) AsInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}
