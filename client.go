package bonusly

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/evergreen-ci/utility"
	"github.com/k0kubun/pp"
	"github.com/mongodb/grip"
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
	catcher := grip.NewBasicCatcher()
	catcher.NewWhen(o.AccessToken == "", "must specify an access token")
	if o.HTTPClient == nil {
		o.HTTPClient = utility.GetDefaultHTTPRetryableClient()
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

func (c *client) CreateBonus(ctx context.Context, opts CreateBonusRequest) (*CreateBonusResponse, error) {
	body, err := c.makeBody(opts)
	if err != nil {
		return nil, errors.Wrap(err, "creating request body")
	}

	r, err := http.NewRequest(http.MethodPost, c.opts.BaseURL+"/bonus", body)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	var result createBonusResponseWrapper
	if err := c.doRequest(ctx, r, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return &result.Result, nil
}

func (c *client) doRequest(ctx context.Context, r *http.Request, result interface{}) error {
	c.addHeaders(r)
	r = r.WithContext(ctx)

	resp, err := c.opts.HTTPClient.Do(r)
	if err != nil {
		pp.Println("making response")
		return errors.Wrap(err, "executing request")
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		pp.Println("reading response")
		return errors.Wrap(err, "reading response body")
	}
	if resp.StatusCode != http.StatusOK {
		pp.Println("non-OK response")
		return c.errorResponse(resp, b)
	}

	if result != nil {
		if err := json.Unmarshal(b, &result); err != nil {
			pp.Println("unmarshalling success result")
			return errors.Wrap(err, "received unexpected response body")
		}

	}

	return nil
}

func (c *client) MyUserInfo(ctx context.Context) (*UserInfoResponse, error) {
	r, err := http.NewRequest(http.MethodGet, c.opts.BaseURL+"/users/me", nil)
	if err != nil {
		pp.Println("NewRequest")
		return nil, errors.Wrap(err, "creating request")
	}
	var result userInfoResponseWrapper
	if err := c.doRequest(ctx, r, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return &result.Result, nil
}

func (c *client) Close(ctx context.Context) error {
	if c.opts.defaultHTTPClient {
		utility.PutHTTPClient(c.opts.HTTPClient)
	}
	return nil
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
		return errors.Wrap(err, statusErr.Error())
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
