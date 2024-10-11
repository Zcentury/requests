package requests

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Zcentury/logger"
	"github.com/Zcentury/requests/method"
	"github.com/Zcentury/requests/params"
	"golang.org/x/net/proxy"
	"io"
	"net"
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

func (r *Requests) Request(m method.Method, url string, args ...any) (*Response, error) {

	argMap, contentType := resolvingArgs(args...)

	data := ""

	if url == "" {
		logger.Error("请传入URL")
		return nil, errors.New("请传入URL")
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
			logger.Error("请传入Body")
			return nil, errors.New("请传入Body")
		}
		// 创建一个新的POST请求
		request, err = http.NewRequest(method.POST.String(), url, strings.NewReader(data))
		//request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
			logger.Error("创建请求失败:%s", err)
			return nil, errors.New("创建请求失败")
		}
	default:
		logger.Error("不支持的请求方式")
		return nil, errors.New("不支持的请求方式")
	}

	if value, ok := argMap["Headers"]; ok {
		request.Header.Add("Content-Type", contentType)
		for k, v := range value.(params.Headers) {
			if k == "Content-Type" {
				request.Header.Set(k, v)
			} else {
				request.Header.Add(k, v)
			}
		}
	}

	if value, ok := argMap["Proxy"]; ok {

		var httpProxy string
		var httpsProxy string
		var socksProxy string

		if v, ook := value.(params.Proxy)["http"]; ook {
			httpProxy = v
		}
		if v, ook := value.(params.Proxy)["https"]; ook {
			httpsProxy = v
		}
		if v, ook := value.(params.Proxy)["socks5"]; ook {
			socksProxy = v
		}

		if httpProxy != "" {
			transport := &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // 跳过TLS证书验证，生产环境中请慎用
				},
			}
			if p, ok := value.(params.Proxy)["http"]; ok {
				proxyURL, err := urlutil.Parse(p)
				if err != nil {
					logger.Error("解析代理地址失败")
					return nil, errors.New("解析代理地址失败")
				}
				transport.Proxy = http.ProxyURL(proxyURL)
			}
			r.client.Transport = transport
		}

		if httpsProxy != "" {
			transport := &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // 跳过TLS证书验证，生产环境中请慎用
				},
			}
			if p, ok := value.(params.Proxy)["http"]; ok {
				proxyURL, err := urlutil.Parse(p)
				if err != nil {
					logger.Error("解析代理地址失败")
					return nil, errors.New("解析代理地址失败")
				}
				transport.Proxy = http.ProxyURL(proxyURL)
			}
			r.client.Transport = transport
		}

		if socksProxy != "" {
			// 创建一个 SOCKS5 代理 Dialer

			dialer, err := proxy.SOCKS5("tcp", socksProxy, nil, proxy.Direct)
			if err != nil {
				return nil, errors.New("创建 SOCKS5 代理失败")
			}

			// 创建一个自定义的 HTTP Transport，并设置 DialContext
			transport := &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return dialer.Dial(network, addr)
				},
			}

			r.client.Transport = transport
		}

	}

	//发送
	response, err := r.client.Do(request)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error("读取响应失败:%s", err)
		return nil, fmt.Errorf("读取响应失败:%s", err)
	}
	bodyBackup := bytes.NewBuffer(body)
	return &Response{
		Response:   response,
		StatusCode: response.StatusCode,
		Text:       string(body),
		Body:       io.NopCloser(bodyBackup),
		Header:     response.Header,
	}, nil
}

// 解析参数
func resolvingArgs(args ...interface{}) (map[string]interface{}, string) {
	result := make(map[string]interface{})
	var contentType string

	for _, arg := range args {

		switch v := reflect.TypeOf(arg); v.Kind() {

		case reflect.String:
			switch v.String() {
			case "params.UrlParams":
				result["UrlParams"] = string(arg.(params.UrlParams))
			case "params.BodyJsonString":
				result["Body"] = string(arg.(params.BodyJsonString))
				contentType = "application/json"
			case "params.BodyString":
				result["Body"] = string(arg.(params.BodyString))
				contentType = "application/x-www-form-urlencoded"
			default:
				logger.Error("未能识别的参数类型")
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
				contentType = "application/json"

			case "params.BodyMap2String":
				values := urlutil.Values{}
				for key, value := range arg.(params.BodyMap2String) {
					values.Add(key, value)
				}
				result["Body"] = values.Encode()
				contentType = "application/x-www-form-urlencoded"

			case "params.Proxy":
				result["Proxy"] = arg.(params.Proxy)

			default:
				logger.Error("未能识别的参数")
			}

		default:
			logger.Error("未能识别的参数")
		}

	}
	return result, contentType
}
