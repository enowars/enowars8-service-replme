package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"replme/types"

	"github.com/gorilla/websocket"
)

type ProxyService struct {
	Target Target
}

type Target struct {
	IP     string
	Port   uint16
	ApiKey string
}

func Proxy(ip string, port uint16, apiKey string) ProxyService {
	return ProxyService{
		Target: Target{
			IP:     ip,
			Port:   port,
			ApiKey: apiKey,
		},
	}
}

func (proxy *ProxyService) SendRegisterRequest(
	request types.RegisterRequest,
	options *types.RequestOptions,
) (*http.Response, *types.RequestError) {
	payload, err := json.Marshal(request)

	if err != nil {
		data, _ := json.Marshal(types.ErrorResponse{
			Error: err.Error(),
		})
		return nil, &types.RequestError{
			Code:        http.StatusBadRequest,
			ContentType: "application/json",
			Data:        data,
		}
	}

	retries := 1
	if options != nil {
		retries = options.Retries
	}

	url := fmt.Sprintf(
		"http://%s:%d/api/%s/auth/register",
		proxy.Target.IP,
		proxy.Target.Port,
		proxy.Target.ApiKey,
	)

	var response *http.Response
	for i := 0; i < retries; i++ {
		response, err = http.Post(
			url,
			"application/json",
			bytes.NewBuffer(payload),
		)
		if err == nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	if err != nil {
		data, _ := json.Marshal(types.ErrorResponse{
			Error: err.Error(),
		})
		return nil, &types.RequestError{
			Code:        http.StatusBadRequest,
			ContentType: "application/json",
			Data:        data,
		}
	}

	if response.StatusCode >= 400 {
		payload, err := io.ReadAll(response.Body)
		if err != nil {
			data, _ := json.Marshal(types.ErrorResponse{
				Error: "Container communication failed",
			})
			return nil, &types.RequestError{
				Code:        http.StatusInternalServerError,
				ContentType: "application/json",
				Data:        data,
			}
		}
		return nil, &types.RequestError{
			Code:        response.StatusCode,
			ContentType: "application/json",
			Data:        payload,
		}
	}

	return response, nil
}

func (proxy *ProxyService) SendCreateTermRequest(
	options types.RequestOptions,
) (*http.Response, *types.RequestError) {

	retries := 1
	if options.Retries > 1 {
		retries = options.Retries
	}

	url := fmt.Sprintf(
		"http://%s:%d/api/%s/term",
		proxy.Target.IP,
		proxy.Target.Port,
		proxy.Target.ApiKey,
	)

	var response *http.Response
	var err error
	for i := 0; i < retries; i++ {
		req, _ := http.NewRequest("POST", url, nil)
		for _, cookie := range options.Cookies {
			req.AddCookie(cookie)
		}
		response, err = http.DefaultClient.Do(req)
		if err == nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	if err != nil {
		data, _ := json.Marshal(types.ErrorResponse{
			Error: err.Error(),
		})
		return nil, &types.RequestError{
			Code:        http.StatusBadRequest,
			ContentType: "application/json",
			Data:        data,
		}
	}

	if response.StatusCode >= 400 {
		payload, err := io.ReadAll(response.Body)
		if err != nil {
			data, _ := json.Marshal(types.ErrorResponse{
				Error: "Container communication failed",
			})
			return nil, &types.RequestError{
				Code:        http.StatusInternalServerError,
				ContentType: "application/json",
				Data:        data,
			}
		}
		return nil, &types.RequestError{
			Code:        response.StatusCode,
			ContentType: "application/json",
			Data:        payload,
		}
	}

	return response, nil
}

func (proxy *ProxyService) SendResizeTermRequest(
	request types.ResizeTermRequest,
	options types.RequestOptions,
) (*http.Response, *types.RequestError) {
	payload, err := json.Marshal(request)
	if err != nil {
		data, _ := json.Marshal(types.ErrorResponse{
			Error: err.Error(),
		})
		return nil, &types.RequestError{
			Code:        http.StatusBadRequest,
			ContentType: "application/json",
			Data:        data,
		}
	}
	retries := 1
	if options.Retries > 1 {
		retries = options.Retries
	}

	url := fmt.Sprintf(
		"http://%s:%d/api/%s/term/size",
		proxy.Target.IP,
		proxy.Target.Port,
		proxy.Target.ApiKey,
	)

	var response *http.Response
	for i := 0; i < retries; i++ {
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range options.Cookies {
			req.AddCookie(cookie)
		}
		response, err = http.DefaultClient.Do(req)
		if err == nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	if err != nil {
		data, _ := json.Marshal(types.ErrorResponse{
			Error: err.Error(),
		})
		return nil, &types.RequestError{
			Code:        http.StatusBadRequest,
			ContentType: "application/json",
			Data:        data,
		}
	}

	if response.StatusCode >= 400 {
		payload, err := io.ReadAll(response.Body)
		if err != nil {
			data, _ := json.Marshal(types.ErrorResponse{
				Error: "Container communication failed",
			})
			return nil, &types.RequestError{
				Code:        http.StatusInternalServerError,
				ContentType: "application/json",
				Data:        data,
			}
		}
		return nil, &types.RequestError{
			Code:        response.StatusCode,
			ContentType: "application/json",
			Data:        payload,
		}
	}

	return response, nil
}

func (proxy *ProxyService) CreateWebsocketPipe(clientConn *websocket.Conn, cookie http.Cookie) error {
	targetURL, err := url.Parse(
		fmt.Sprintf(
			"ws://%s:%d/api/%s/term",
			proxy.Target.IP,
			proxy.Target.Port,
			proxy.Target.ApiKey,
		),
	)

	if err != nil {
		return err
	}

	targetConn, _, err := websocket.DefaultDialer.Dial(
		targetURL.String(),
		http.Header{
			"Cookie": []string{
				cookie.String(),
			},
		},
	)

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
