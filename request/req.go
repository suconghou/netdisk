package request

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"netdisk/util"
	"time"
)

var (
	dialer = &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
)

type DNSMapper struct {
	mappings map[string]string
}

// DialContext 是自定义的拨号函数，用于替换域名解析逻辑
func (m *DNSMapper) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	// 如果域名在映射表中，替换为指定 IP
	if ip, ok := m.mappings[host]; ok {
		addr = net.JoinHostPort(ip, port)
	}
	// 使用默认的 Dialer 建立连接, addr必须包含端口号
	return dialer.DialContext(ctx, network, addr)
}

// host 指定固定的DNS解析值，值为一个IP或域名，不能包含端口，请求将发送给这个IP地址
func Request(ctx context.Context, urlStr string, method string, reqHeader http.Header, timeout int64, body io.Reader, host string) (*http.Response, bool, error) {
	t := time.Second * time.Duration(timeout)
	ctx2, cancel := context.WithTimeout(ctx, t)
	req, err := http.NewRequestWithContext(ctx2, method, urlStr, body)
	if err != nil {
		cancel()
		return nil, false, err
	}
	req.Header = reqHeader
	var (
		resp *http.Response
	)
	if host != "" {
		mapper := &DNSMapper{
			mappings: map[string]string{
				req.URL.Hostname(): host, // 指定IP或者域名
			},
		}
		cli := http.Client{
			Timeout: t,
			Transport: &http.Transport{
				Proxy:                 http.ProxyFromEnvironment,
				DialContext:           mapper.DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			},
		}
		resp, err = cli.Do(req)
	} else {
		resp, err = http.DefaultClient.Do(req)
	}
	if err != nil {
		cancel()
		return resp, false, err
	}
	resp.Body = util.NewOnCloseReadCloser(resp.Body, func() error {
		cancel()
		return nil
	})
	statusOk := resp.StatusCode/100 == 2
	return resp, statusOk, nil
}
