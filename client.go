package bonusly

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// ClientOptions represent options to initialize a Bonusly client authenticated
// with a particular user's access token.
type ClientOptions struct {
	AccessToken       string
	BaseURL           string
	HTTPClient        *http.Client
	defaultHTTPClient bool
}

// Validate checks that all the required fields are set and sets defaults where
// possible.
func (o *ClientOptions) Validate() error {
	catcher := newBasicCatcher()
	catcher.NewWhen(o.AccessToken == "", "must specify an access token")
	if o.HTTPClient == nil {
		o.HTTPClient = getDefaultHTTPRetryableClient()
		o.defaultHTTPClient = true
	}
	if o.BaseURL == "" {
		o.BaseURL = productionBaseURL
	}
	o.BaseURL = strings.TrimSuffix(o.BaseURL, "/")
	return catcher.Resolve()
}

type client struct {
	opts ClientOptions
}

// NewClient returns a client to interact with the Bonusly API.
func NewClient(opts ClientOptions) (Client, error) {
	if err := opts.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid options")
	}
	return &client{
		opts: opts,
	}, nil
}

func (c *client) CreateBonus(ctx context.Context, opts CreateBonusRequest) (*BonusResponse, error) {
	body, err := c.makeBody(opts)
	if err != nil {
		return nil, errors.Wrap(err, "creating request body")
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, c.urlRoute("/bonuses"), body)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	var result bonusResponseWrapper
	if err := c.doRequest(r, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return &result.Result, nil
}

func (c *client) GetBonus(ctx context.Context, id string) (*BonusResponse, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, c.urlRoute("/bonuses", id), nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	var result bonusResponseWrapper
	if err := c.doRequest(r, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return &result.Result, nil
}

func (c *client) ListBonuses(ctx context.Context, req ListBonusesRequest) ([]BonusResponse, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, c.urlRoute("/bonuses"), nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}
	q := r.URL.Query()
	for k, v := range req.QueryMap() {
		q.Set(k, v)
	}
	r.URL.RawQuery = q.Encode()

	var result bonusesResponseWrapper
	if err := c.doRequest(r, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result.Result, nil
}

func (c *client) UpdateBonus(ctx context.Context, id, reason string) (*BonusResponse, error) {
	payload := CreateBonusRequest{
		Reason: reason,
	}
	body, err := c.makeBody(payload)
	if err != nil {
		return nil, errors.Wrap(err, "creating request body")
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPut, c.urlRoute("/bonuses", id), body)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	var result bonusResponseWrapper
	if err := c.doRequest(r, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return &result.Result, nil
}

func (c *client) DeleteBonus(ctx context.Context, id string) error {
	r, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.urlRoute("/bonuses", id), nil)
	if err != nil {
		return errors.Wrap(err, "creating request")
	}

	if err := c.doRequest(r, nil); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *client) ListRewards(ctx context.Context, req ListRewardsRequest) ([]RewardsResponse, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, c.urlRoute("/rewards"), nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}
	q := r.URL.Query()
	for k, v := range req.QueryMap() {
		q.Set(k, v)
	}
	r.URL.RawQuery = q.Encode()

	var result rewardsResponseWrapper
	if err := c.doRequest(r, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result.Result, nil
}

func (c *client) MyUserInfo(ctx context.Context) (*UserInfoResponse, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, c.urlRoute("/users/me"), nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}
	var result userInfoResponseWrapper
	if err := c.doRequest(r, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return &result.Result, nil
}

func (c *client) Close(_ context.Context) error {
	if c.opts.defaultHTTPClient {
		putHTTPClient(c.opts.HTTPClient)
	}
	return nil
}

func (c *client) doRequest(r *http.Request, result interface{}) error {
	c.addHeaders(r)

	resp, err := c.opts.HTTPClient.Do(r)
	if err != nil {
		return errors.Wrap(err, "executing request")
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "reading response body")
	}
	if resp.StatusCode != http.StatusOK {
		return c.errorResponse(resp, b)
	}

	if result != nil {
		if err := json.Unmarshal(b, &result); err != nil {
			return errors.Wrap(err, "received unexpected response body")
		}
	}

	return nil
}

func (c *client) urlRoute(parts ...string) string {
	baseURL := strings.TrimSuffix(c.opts.BaseURL, "/")
	if len(parts) == 0 {
		return baseURL
	}
	parts[0] = strings.TrimPrefix(parts[0], "/")
	return fmt.Sprintf("%s/%s", baseURL, path.Join(parts...))
}

func (c *client) makeBody(payload interface{}) (io.Reader, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrapf(err, "marshalling payload of type %T to JSON", payload)
	}
	return bufio.NewReader(bytes.NewBuffer(b)), nil
}

func (c *client) errorResponse(resp *http.Response, body []byte) error {
	var errResp CommonResponse
	statusErr := errors.Errorf("status %s", resp.Status)
	if err := json.Unmarshal(body, &errResp); err != nil {
		return errors.Wrap(errors.New(string(body)), statusErr.Error())
	}
	if errResp.Message != nil {
		return errors.Wrap(errors.New(*errResp.Message), statusErr.Error())
	}
	if !fromBoolPtr(errResp.Success) {
		return errors.Wrap(errors.New("request unsuccessful for unknown reason"), statusErr.Error())
	}
	return errors.WithStack(statusErr)
}

func (c *client) addHeaders(r *http.Request) {
	r.Header.Add("Content-Type", contentType)
	r.Header.Add("Authorization", "Bearer "+c.opts.AccessToken)
}
