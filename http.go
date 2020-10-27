package bonusly

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/PuerkitoBio/rehttp"
)

// Source: github.com/evergreen-ci/utility

const httpClientTimeout = 5 * time.Minute

var httpClientPool *sync.Pool

func init() {
	initHTTPPool()
}

func initHTTPPool() {
	httpClientPool = &sync.Pool{
		New: func() interface{} { return newBaseConfiguredHTTPClient() },
	}
}

func newBaseConfiguredHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   httpClientTimeout,
		Transport: newConfiguredBaseTransport(),
	}
}

func newConfiguredBaseTransport() *http.Transport {
	return &http.Transport{
		TLSClientConfig:     &tls.Config{},
		Proxy:               http.ProxyFromEnvironment,
		DisableCompression:  false,
		DisableKeepAlives:   true,
		IdleConnTimeout:     20 * time.Second,
		MaxIdleConnsPerHost: 10,
		MaxIdleConns:        50,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 0,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

}

// getHTTPClient produces default HTTP client from the pool,
// constructing a new client if needed. Always pair calls to
// getHTTPClient with defered calls to putHTTPClient.
func getHTTPClient() *http.Client { return httpClientPool.Get().(*http.Client) }

// putHTTPClient returns the client to the pool, automatically
// reconfiguring the transport.
func putHTTPClient(c *http.Client) {
	c.Timeout = httpClientTimeout

	switch transport := c.Transport.(type) {
	case *http.Transport:
		transport.TLSClientConfig.InsecureSkipVerify = false
		c.Transport = transport
	case *rehttp.Transport:
		c.Transport = transport.RoundTripper
		putHTTPClient(c)
		return
	default:
		c.Transport = newConfiguredBaseTransport()
	}

	httpClientPool.Put(c)
}

// httpRetryConfiguration makes it possible to configure the retry
// semantics for retryable clients. In most cases, construct this
// object using the newDefaultHTTPRetryConf, which provides reasonable
// defaults.
type httpRetryConfiguration struct {
	MaxDelay        time.Duration
	BaseDelay       time.Duration
	MaxRetries      int
	TemporaryErrors bool
	Methods         []string
	Statuses        []int
	Errors          []error
	ErrorStrings    []string
}

// newDefaultHTTPRetryConf constructs a httpRetryConfiguration object
// with reasonable defaults.
func newDefaultHTTPRetryConf() httpRetryConfiguration {
	return httpRetryConfiguration{
		MaxRetries:      50,
		TemporaryErrors: true,
		MaxDelay:        5 * time.Second,
		BaseDelay:       50 * time.Millisecond,
		Methods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
		},
		Statuses: []int{
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
			http.StatusInsufficientStorage,
			http.StatusConflict,
			http.StatusRequestTimeout,
			http.StatusPreconditionFailed,
			http.StatusExpectationFailed,
		},
	}
}

// getHTTPRetryableClient produces an HTTP client that automatically
// retries failed requests according to the configured
// parameters. Couple calls to getHTTPRetryableClient, with defered
// calls to putHTTPClient.
func getHTTPRetryableClient(conf httpRetryConfiguration) *http.Client {
	client := getHTTPClient()

	statusRetries := []rehttp.RetryFn{}
	if len(conf.Statuses) > 0 {
		statusRetries = append(statusRetries, rehttp.RetryStatuses(conf.Statuses...))
	} else {
		conf.TemporaryErrors = true
	}

	if conf.TemporaryErrors {
		statusRetries = append(statusRetries, rehttp.RetryTemporaryErr())
	}

	if len(conf.Errors) > 0 {
		statusRetries = append(statusRetries, rehttp.RetryIsErr(func(err error) bool {
			for _, errToCheck := range conf.Errors {
				if err == errToCheck {
					return true
				}
			}
			return false
		}))
	}

	if len(conf.ErrorStrings) > 0 {
		statusRetries = append(statusRetries, rehttp.RetryIsErr(func(err error) bool {
			for _, errToCheck := range conf.ErrorStrings {
				if err.Error() == errToCheck {
					return true
				}
			}
			return false
		}))
	}

	retryFns := []rehttp.RetryFn{rehttp.RetryAny(statusRetries...)}

	if len(conf.Methods) > 0 {
		retryFns = append(retryFns, rehttp.RetryHTTPMethods(conf.Methods...))
	}

	if conf.MaxRetries > 0 {
		retryFns = append(retryFns, rehttp.RetryMaxRetries(conf.MaxRetries))
	}

	client.Transport = rehttp.NewTransport(client.Transport,
		rehttp.RetryAll(retryFns...),
		rehttp.ExpJitterDelay(conf.BaseDelay, conf.MaxDelay))

	return client
}

// getDefaultHTTPRetryableClient provides a retryable client with the default
// settings. Couple calls to getDefaultHTTPRetryableClient, with defered calls
// to putHTTPClient.
func getDefaultHTTPRetryableClient() *http.Client {
	return getHTTPRetryableClient(newDefaultHTTPRetryConf())
}

// httpRetryFunction makes it possible to write customizable retry
// logic. Returning true if the request should be retried again and
// false otherwise.
type httpRetryFunction func(index int, req *http.Request, resp *http.Response, err error) bool

// httpDelayFunction makes it possible to write customizable retry
// backoff logic, by allowing you to evaluate the previous request and
// response and return the duration to wait before the next request.
type httpDelayFunction func(index int, req *http.Request, resp *http.Response, err error) time.Duration

func makeRetryFn(in httpRetryFunction) rehttp.RetryFn {
	return func(attempt rehttp.Attempt) bool {
		return in(attempt.Index, attempt.Request, attempt.Response, attempt.Error)
	}
}

func makeDelayFn(in httpDelayFunction) rehttp.DelayFn {
	return func(attempt rehttp.Attempt) time.Duration {
		return in(attempt.Index, attempt.Request, attempt.Response, attempt.Error)
	}
}

// getCustomHTTPRetryableClient allows you to generate an HTTP client
// that automatically retries failed request based on the provided
// custom logic.
func getCustomHTTPRetryableClient(retry httpRetryFunction, delay httpDelayFunction) *http.Client {
	client := getHTTPClient()
	client.Transport = rehttp.NewTransport(client.Transport, makeRetryFn(retry), makeDelayFn(delay))
	return client
}
