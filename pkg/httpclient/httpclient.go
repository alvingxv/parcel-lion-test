package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"go.elastic.co/apm/module/apmhttp"
)

type HttpClient struct {
	client      *http.Client
	fallbackMsg map[string]string
	source      hystrixSource
}

type hystrixSource struct {
	command                string
	timeout                time.Duration
	maxConcurrentRequests  int
	errorPercentThreshold  int
	sleepWindow            int
	requestVolumeThreshold int
	fallbackMsg            string
}

const (
	defaultHystrixRetryCount      = 0
	defaultHTTPTimeout            = 30 * time.Second
	defaultHystrixTimeout         = 30 * time.Second
	defaultMaxConcurrentRequests  = 100
	defaultErrorPercentThreshold  = 25
	defaultSleepWindow            = 10
	defaultRequestVolumeThreshold = 10
	defaultFallbackMsg            = "Timeout External (State: Open)"
	defaultCommand                = "http-call"
)

type Option func(*HttpClient)

var Client *HttpClient

func Init() {
	Client = &HttpClient{
		client: apmhttp.WrapClient(&http.Client{
			Timeout: 20 * time.Second,
		}),
		fallbackMsg: make(map[string]string),
	}
}

func (c *HttpClient) NewCbSource(opts ...Option) {

	c.source.errorPercentThreshold = defaultErrorPercentThreshold
	c.source.maxConcurrentRequests = defaultMaxConcurrentRequests
	c.source.sleepWindow = defaultSleepWindow
	c.source.requestVolumeThreshold = defaultRequestVolumeThreshold
	c.source.timeout = defaultHystrixTimeout
	c.source.fallbackMsg = defaultFallbackMsg
	c.source.command = defaultCommand

	for _, opt := range opts {
		opt(c)
	}

	hystrix.ConfigureCommand(c.source.command, hystrix.CommandConfig{
		Timeout:                int(c.source.timeout),
		MaxConcurrentRequests:  c.source.maxConcurrentRequests,
		ErrorPercentThreshold:  c.source.errorPercentThreshold,
		RequestVolumeThreshold: c.source.requestVolumeThreshold,
		SleepWindow:            c.source.sleepWindow,
	})

	c.fallbackMsg[c.source.command] = c.source.fallbackMsg
}

func (c *HttpClient) CbWithTimeout(timeout time.Duration) Option {
	return func(c *HttpClient) {
		c.source.timeout = timeout
	}
}
func (c *HttpClient) CbWithCommand(command string) Option {
	return func(c *HttpClient) {
		c.source.command = command
	}
}
func (c *HttpClient) CbWithMaxConcurrentRequests(maxConcurrentRequests int) Option {
	return func(c *HttpClient) {
		c.source.maxConcurrentRequests = maxConcurrentRequests
	}
}
func (c *HttpClient) CbWithErrorPercentThreshold(errorPercentThreshold int) Option {
	return func(c *HttpClient) {
		c.source.errorPercentThreshold = errorPercentThreshold
	}
}
func (c *HttpClient) CbWithRequestVolumeThreshold(requestVolumeThreshold int) Option {
	return func(c *HttpClient) {
		c.source.requestVolumeThreshold = requestVolumeThreshold
	}
}
func (c *HttpClient) CbWithSleepWindow(sleepWindow int) Option {
	return func(c *HttpClient) {
		c.source.sleepWindow = sleepWindow
	}
}
func (c *HttpClient) CbWithFallbackMsg(fallbackMsg string) Option {
	return func(c *HttpClient) {
		c.source.fallbackMsg = fallbackMsg
	}
}

func (c *HttpClient) Call(ctx context.Context, requestBody map[string]interface{}, header http.Header, endpoint string, source string) (context.Context, []byte, http.Header, error) {
	var responseCtx context.Context
	var responseBody []byte
	var responseHeader http.Header
	var responseErr error

	err := hystrix.Do(source,
		func() error {
			var err error
			responseCtx, responseBody, responseHeader, err = c.makeHttpCall(ctx, requestBody, header, endpoint)
			if err != nil {
				// log.LogDebug(fmt.Sprint("Main call failed for %s: %v", source, err))
			}
			return err
		},
		func(err error) error {
			// log.LogDebug(fmt.Sprint("Fallback triggered for %s. Error: %v", source, err))
			fmt.Println("fallback")

			responseCtx = ctx
			responseBody = nil
			responseHeader = make(http.Header)
			responseHeader.Set("Content-Type", "application/json")
			responseErr = fmt.Errorf(c.fallbackMsg[source])

			return nil
		})

	if err != nil {
		// log.LogDebug(fmt.Sprint("Complete failure in Call to %s: %v", source, err))

		return ctx, nil, nil, fmt.Errorf("service completely unavailable: %v", err)
	}

	return responseCtx, responseBody, responseHeader, responseErr
}

func (c *HttpClient) makeHttpCall(ctx context.Context, requestBody map[string]interface{}, header http.Header, endpoint string) (context.Context, []byte, http.Header, error) {
	jsonRequest, _ := json.Marshal(requestBody)

	payload := bytes.NewReader(jsonRequest)

	request, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		return ctx, nil, nil, err
	}

	request.Header = header

	response, err := c.client.Do(request.WithContext(ctx))
	if err != nil {
		return ctx, nil, nil, err
	}

	defer response.Body.Close()

	responseByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ctx, nil, nil, err
	}

	// check for valid json response
	var js map[string]interface{}
	err = json.Unmarshal(responseByte, &js)
	if err != nil {
		//log.LogWarn("invalid json", "invalid json")
		return ctx, nil, nil, err
	}

	return ctx, responseByte, response.Header, nil
}
