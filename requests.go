package requests

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/Zcentury/gologger"
	"github.com/Zcentury/requests/method"
	"github.com/Zcentury/requests/params"
	"io"
	"net/http"
	urlutil "net/url"
	"reflect"
	"strings"
)

type Requests struct {
	client *http.Client
}

func NewRequests(client *http.Client) *Requests {
	return &Requests{
		client: client,
	}
}

func (r *Requests) SetProxy(proxy string) error {
	proxyURL, err := urlutil.Parse(proxy)
	if err != nil {
		gologger.Error().Msg("解析代理地址失败")
		return errors.New("解析代理地址失败")
	}

	//创建一个Transport并设置代理
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 跳过TLS证书验证，生产环境中请慎用
		},
	}

	r.client.Transport = transport

	return nil
}

func (r *Requests) Request(m method.Method, url string, args ...any) *Response {

	argMap := resolvingArgs(args...)

	data := ""

	if url == "" {
		gologger.Error().Msg("请传入URL")
	}

	if value, ok := argMap["UrlParams"]; ok {
		url += "?" + value.(string)
	}

	var request *http.Request
	var err error
	switch m {
	case method.GET:
		request, err = http.NewRequest(method.GET.String(), url, strings.NewReader(data))
	case method.POST:
		if value, ok := argMap["Body"]; ok {
			data = value.(string)
		} else {
			gologger.Error().Msg("请传入Body")
			return nil
		}
		// 创建一个新的POST请求
		request, err = http.NewRequest(method.POST.String(), url, strings.NewReader(data))
		//request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
			gologger.Error().Msgf("创建请求失败:%s", err)
			return nil
		}
	default:
		gologger.Error().Msg("不支持的请求方式")
		return nil
	}

	if value, ok := argMap["Headers"]; ok {
		for k, v := range value.(params.Headers) {
			request.Header.Add(k, v)
		}
	}

	if value, ok := argMap["Proxy"]; ok {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 跳过TLS证书验证，生产环境中请慎用
			},
		}
		if proxy, ok := value.(params.Proxy)["http"]; ok {
			proxyURL, err := urlutil.Parse(proxy)
			if err != nil {
				gologger.Error().Msg("解析代理地址失败")
				return nil
			}
			transport.Proxy = http.ProxyURL(proxyURL)
		}
		//if proxy, ok := value.(params.Proxy)["http"]; ok {
		//
		//}
		r.client.Transport = transport
	}

	//发送
	response, err := r.client.Do(request)
	if err != nil {
		gologger.Error().Msg(err.Error())
		return nil
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		gologger.Error().Msgf("读取响应失败:%s", err)
		return nil
	}

	return &Response{
		Response: response,
		Text:     string(body),
	}
}

func (r *Requests) Get(url string, args ...interface{}) *Response {
	return r.Request(method.GET, url, args...)
}

func (r *Requests) Post(url string, args ...interface{}) *Response {
	return r.Request(method.POST, url, args...)
}

// 解析参数
func resolvingArgs(args ...interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for _, arg := range args {

		switch v := reflect.TypeOf(arg); v.Kind() {

		case reflect.String:
			switch v.String() {
			case "params.UrlParams":
				result["UrlParams"] = string(arg.(params.UrlParams))
			case "params.BodyJsonString":
				result["Body"] = string(arg.(params.BodyJsonString))
			case "params.BodyString":
				result["Body"] = string(arg.(params.BodyString))
			default:
				gologger.Error().Msg("未能识别的参数类型")
			}

		case reflect.Map:

			switch v.String() {
			case "params.Headers":
				result["Headers"] = arg.(params.Headers)

			case "params.UrlMap2Params":
				values := urlutil.Values{}
				for key, value := range arg.(params.UrlMap2Params) {
					values.Add(key, value)
				}
				result["UrlParams"] = values.Encode()

			case "params.BodyMap2Json":
				if jsonData, err := json.Marshal(arg.(params.BodyMap2Json)); err == nil && result["Body"] != "" {
					result["Body"] = string(jsonData)
				}

			case "params.BodyMap2String":
				values := urlutil.Values{}
				for key, value := range arg.(params.BodyMap2String) {
					values.Add(key, value)
				}
				result["Body"] = values.Encode()

			case "params.Proxy":
				result["Proxy"] = arg.(params.Proxy)

			default:
				gologger.Error().Msg("未能识别的参数")
			}

		default:
			gologger.Error().Msg("未能识别的参数")
		}

	}
	return result
}
