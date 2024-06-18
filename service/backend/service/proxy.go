package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"replme/types"
	"replme/util"

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
		util.SLogger.Debugf("[UN:%s..] -> [%s:%d]  Sending register request (try %d)", request.Username[:5], proxy.Target.IP, proxy.Target.Port, i+1)
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
		util.SLogger.Errorf("[UN:%s..] -> [%s:%d] Register request failed: %s", request.Username[:5], proxy.Target.IP, proxy.Target.Port, err.Error())
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
			util.SLogger.Warnf("[UN:%s..] -> [%s:%d] Register request failed: [%d] %s", request.Username[:5], proxy.Target.IP, proxy.Target.Port, response.StatusCode, err.Error())
			return nil, &types.RequestError{
				Code:        http.StatusInternalServerError,
				ContentType: "application/json",
				Data:        data,
			}
		}
		util.SLogger.Warnf("[UN:%s..] -> [%s:%d] Register request failed: [%d] %s", request.Username[:5], proxy.Target.IP, proxy.Target.Port, response.StatusCode, string(payload))
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

	util.SLogger.Debugf("[..] -> [%s:%d] Dialing websocket", proxy.Target.IP, proxy.Target.Port)
	targetConn, _, err := websocket.DefaultDialer.Dial(
		targetURL.String(),
		http.Header{
			"Cookie": []string{
				cookie.String(),
			},
		},
	)

	if err != nil {
		util.SLogger.Errorf("[..] -> [%s:%d] Dialing failed: ", proxy.Target.IP, proxy.Target.Port, err.Error())
		return err
	}

	defer targetConn.Close()
	defer func() {
		util.SLogger.Debugf("[..] -> [%s:%d] Closing websocket pipe", proxy.Target.IP, proxy.Target.Port)
	}()

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
		util.SLogger.Errorf("[..] -> [%s:%d] Websocket pipe failed: ", proxy.Target.IP, proxy.Target.Port, err.Error())
	}

	return nil
}
