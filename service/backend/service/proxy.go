package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"cafedodo/types"

	"github.com/gorilla/websocket"
)

type ProxyService struct {
	Target Target
}

type Target struct {
	IP   string
	Port uint16
}

func Proxy(ip string, port uint16) ProxyService {
	return ProxyService{
		Target: Target{
			IP:   ip,
			Port: port,
		},
	}
}

func (proxy *ProxyService) SendUserCreateRequest(
	request types.LoginRequest,
) (*http.Response, error) {
	payload, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	return http.Post(
		fmt.Sprintf(
			"http://%s:%d/user/create",
			proxy.Target.IP,
			proxy.Target.Port,
		),
		"application/json",
		bytes.NewBuffer(payload),
	)
}

func (proxy *ProxyService) SendUserLoginRequest(
	request types.LoginRequest,
) (*http.Response, error) {
	payload, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	return http.Post(
		fmt.Sprintf(
			"http://%s:%d/user/login",
			proxy.Target.IP,
			proxy.Target.Port,
		),
		"application/json",
		bytes.NewBuffer(payload),
	)
}

func (proxy *ProxyService) SendCreateTermProcessRequest(
	request types.CreateTermProcessRequest,
) (*http.Response, error) {
	url := fmt.Sprintf(
		"http://%s:%d/terminals?rows=%d&cols=%d",
		proxy.Target.IP,
		proxy.Target.Port,
		request.Rows,
		request.Columns,
	)
	req, _ := http.NewRequest("POST", url, nil)

	return http.DefaultClient.Do(req)
}

func (proxy *ProxyService) SendUpdateTermSizeRequest(
	pid string, request types.UpdateTermSizeRequest,
) (*http.Response, error) {
	url := fmt.Sprintf(
		"http://%s:%d/terminals/%s/size?rows=%d&cols=%d",
		proxy.Target.IP,
		proxy.Target.Port,
		pid,
		request.Rows,
		request.Columns,
	)
	req, _ := http.NewRequest("POST", url, nil)

	return http.DefaultClient.Do(req)
}

func (proxy *ProxyService) CreateWebsocketPipe(clientConn *websocket.Conn, pid string) error {
	targetURL, err := url.Parse(
		fmt.Sprintf(
			"ws://%s:%d/terminals/%s",
			proxy.Target.IP,
			proxy.Target.Port,
			pid,
		),
	)

	if err != nil {
		return err
	}

	targetConn, _, err := websocket.DefaultDialer.Dial(targetURL.String(), nil)

	if err != nil {
		return err
	}

	defer targetConn.Close()

	errc := make(chan error, 2)

	go func() {
		_, err := io.Copy(clientConn.UnderlyingConn(), targetConn.UnderlyingConn())
		errc <- err
	}()
	go func() {
		_, err := io.Copy(targetConn.UnderlyingConn(), clientConn.UnderlyingConn())
		errc <- err
	}()

	if err := <-errc; err != nil {
		log.Printf("WebSocket proxy error: %v", err)
	}

	return nil
}
