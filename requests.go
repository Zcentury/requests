package requests

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/Zcentury/requests/method"
	"github.com/Zcentury/requests/params"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	urlutil "net/url"
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
	request, err := r.generateRequest(m, url, args...)
	if err != nil {
		return nil, err
	}

	//发送
	response, err := r.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	header := response.Header

	var reader io.ReadCloser
	switch header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
	case "deflate":
		reader, err = zlib.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
	default:
		reader = response.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败:%s", err)
	}

	bodyBackup := bytes.NewBuffer(body)
	return &Response{
		Response:   response,
		StatusCode: response.StatusCode,
		Text:       string(body),
		Body:       io.NopCloser(bodyBackup),
		Header:     header,
	}, nil
}

func (r *Requests) generateRequest(m method.Method, url string, args ...any) (*http.Request, error) {
	var err error

	argMap, contentType, err := resolvingArgs(args...)
	if err != nil {
		return nil, err
	}

	data := ""

	if url == "" {
		return nil, ErrNoUrl
	}

	if value, ok := argMap["UrlParams"]; ok {
		url += "?" + value.(string)
	}

	var request *http.Request

	switch m {
	case method.GET:
		request, err = http.NewRequest(method.GET.String(), url, strings.NewReader(data))
	case method.POST:
		if value, ok := argMap["Body"]; ok {
			data = value.(string)
		} else {
			return nil, ErrNoBody
		}
		// 创建一个新的POST请求
		request, err = http.NewRequest(method.POST.String(), url, strings.NewReader(data))
		//request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
			return nil, fmt.Errorf("创建请求失败:%s", err)
		}
	default:
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
	return request, nil
}
