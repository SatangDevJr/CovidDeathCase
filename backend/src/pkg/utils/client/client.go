package client

import (
	"bytes"
	"context"
	requestHeader "covid/src/api/requestheader"
	"covid/src/pkg/utils/logger"
	"net/http"
	"strings"
)

type Service struct {
	Service UseCase
	Logs    logger.Logger
}

type UseCase interface {
	Request(options HttpOptions) (*http.Response, error)
}

type ServiceParam struct {
	Logs logger.Logger
}

func NewService(logs logger.Logger) *Service {
	return &Service{
		Logs: logs,
	}
}

type HttpOptions struct {
	Context context.Context
	URL     string
	Method  string
	Headers map[string]string
	Queries map[string]string
	Body    *bytes.Buffer
}

type FileInfo struct {
	Key      string
	File     []byte
	FileName string
}

func (service *Service) Request(options HttpOptions) (*http.Response, error) {

	if options.Queries != nil {
		var queryParam []string
		for key, value := range options.Queries {
			queryParam = append(queryParam, key+"="+value)
		}

		pathQuery := strings.Join(queryParam, "&")
		options.URL += "?" + pathQuery
	}

	var request *http.Request
	var err error

	if options.Method == "GET" {
		request, err = http.NewRequest(options.Method, options.URL, nil)
	} else {
		request, err = http.NewRequest(options.Method, options.URL, options.Body)
	}
	if err != nil {
		return nil, err
	}

	if options.Headers[requestHeader.Authorization] != "" {
		request.Header.Set(requestHeader.Authorization, options.Headers[requestHeader.Authorization])
	}
	if options.Headers[requestHeader.ContentType] != "" {
		request.Header.Set(requestHeader.ContentType, options.Headers[requestHeader.ContentType])
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	return response, nil
}