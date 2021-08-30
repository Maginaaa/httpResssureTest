package model

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	HttpOk         = 200 // 请求成功
	RequestTimeout = 506 // 请求超时
	RequestErr     = 509 // 请求错误
	ParseError     = 510 // 解析错误

	FormTypeHttp      = "http"
	TimeOut  = 30 * time.Second

)

// 验证方法
type VerifyHttp func(request *Request, response *http.Response) (code int, isSucceed bool)


// 请求结果
type Request struct {
	Url     string            // Url
	Form    string            // http/webSocket/tcp
	Method  string            // 方法 GET/POST/PUT
	Headers map[string]string // Headers
	Body    string            // body
	Verify  string            // 验证的方法
	Timeout time.Duration     // 请求超时时间
	Debug   bool              // 是否开启Debug模式
}

func (r *Request) GetBody() (body io.Reader) {
	body = strings.NewReader(r.Body)

	return
}

func (r *Request) GetDebug() bool {

	return r.Debug
}


// NewRequest
// url 压测的url
// verify 验证方法 在server/verify中 http 支持:statusCode、json webSocket支持:json
// timeout 请求超时时间
// debug 是否开启debug
// path curl文件路径 http接口压测，自定义参数设置
func NewRequest(url string, reqHeaders []string, reqBody string) (request *Request, err error) {

	var (
		method  = "GET"
		headers = make(map[string]string)
		body    string
	)

	// body取值
	if reqBody != "" {
		method = "POST"
		body = reqBody
		headers["Content-Type"] = "application/x-www-form-urlencoded; charset=utf-8"
	}

	// 获取请求头
	for _, v := range reqHeaders {
		getHeaderValue(v, headers)
	}

	// url拼接
	form := FormTypeHttp
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = fmt.Sprintf("http://%s", url)
	}

	request = &Request{
		Url:     url,
		Form:    form,
		Method:  strings.ToUpper(method),
		Headers: headers,
		Body:    body,
		Timeout: TimeOut,
	}

	return

}

func getHeaderValue(v string, headers map[string]string) {
	index := strings.Index(v, ":")
	if index < 0 {
		return
	}

	vIndex := index + 1
	if len(v) >= vIndex {
		value := strings.TrimPrefix(v[vIndex:], " ")

		if _, ok := headers[v[:index]]; ok {
			headers[v[:index]] = fmt.Sprintf("%s; %s", headers[v[:index]], value)
		} else {
			headers[v[:index]] = value
		}
	}
}

// 打印
func (r *Request) Print() {
	if r == nil {

		return
	}

	result := fmt.Sprintf("request:\n form:%s \n url:%s \n method:%s \n headers:%v \n", r.Form, r.Url, r.Method, r.Headers)
	result = fmt.Sprintf("%s data:%v \n", result, r.Body)
	result = fmt.Sprintf("%s \n timeout:%s \n debug:%v \n", result, r.Timeout, r.Debug)
	fmt.Println(result)

	return
}


// 请求结果
type RequestResults struct {
	Id            string // 消息Id
	ChanId        uint64 // 消息Id
	Time          uint64 // 请求时间 纳秒
	IsSucceed     bool   // 是否请求成功
	ErrCode       int    // 错误码
	ReceivedBytes int64
}

func (r *RequestResults) SetId(chanId uint64, number uint64) {
	id := fmt.Sprintf("%d_%d", chanId, number)

	r.Id = id
	r.ChanId = chanId
}