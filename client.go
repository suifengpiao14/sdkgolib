package sdkgolib

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/suifengpiao14/logchan/v2"
	"github.com/suifengpiao14/torm"
)

var (
	API_NOT_FOUND = errors.Errorf("not found client")
)

type DefaultImplementClientOutput struct{}

type OutI interface {
	Error() (err error)
}

func (c DefaultImplementClientOutput) Error() (err error) {
	return nil
}

type ClientInterface interface {
	GetRoute() (method string, path string)
	Init()
	Request(ctx context.Context) (err error)
	GetDescription() (title string, description string)
	GetName() (domain string, name string)
	GetOutRef() (outRef OutI)
	RequestHandler(ctx context.Context, input []byte) (out []byte, err error)
	GetSDKConfig() (sdkConfig Config)
}

type DefaultImplementPartClientFuncs struct {
}

func (e *DefaultImplementPartClientFuncs) Init() {
}

// RequestFn 封装http请求数据格式
type RequestFn func(ctx context.Context, req *http.Request) (out []byte, err error)

// RestyRequestFn 通用请求方法
func RestyRequestFn(ctx context.Context, req *http.Request) (out []byte, err error) {
	r := resty.New().R()
	urlstr := req.URL.String()
	r.Header = req.Header
	r.FormData = req.Form
	r.RawRequest = req
	if req.Body != nil {
		var body io.ReadCloser
		body, err = req.GetBody()
		if err == nil && body != nil {
			defer body.Close()
			b, err := io.ReadAll(body)
			if err != nil {
				return nil, err
			}
			r.SetBody(b)
		}
	}

	logInfo := &torm.LogInfoHttp{
		GetRequest: func() *http.Request { return r.RawRequest },
	}
	defer func() {
		logchan.SendLogInfo(logInfo)
	}()
	res, err := r.Execute(strings.ToUpper(req.Method), urlstr)
	if err != nil {
		return nil, err
	}

	responseBody := res.Body()
	if res.StatusCode() != http.StatusOK {
		err = errors.Errorf("http status:%d,body:%s", res.StatusCode(), string(responseBody))
		return nil, err
	}
	logInfo.ResponseBody = string(responseBody)
	logInfo.Response = res.RawResponse
	return responseBody, nil

}
