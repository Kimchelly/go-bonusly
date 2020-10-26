package bonusly

import (
	"context"
	"os"
	"testing"

	"github.com/evergreen-ci/utility"
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	accessToken := os.Getenv("BONUSLY_TOKEN")
	require.NotEmpty(t, accessToken, "missing access token")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpClient := utility.GetHTTPClient()
	defer utility.PutHTTPClient(httpClient)

	c, err := NewClient(ClientOptions{
		AccessToken: accessToken,
		HTTPClient:  httpClient,
	})
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, c.Close(ctx))
	}()

	info, err := c.MyUserInfo(ctx)
	require.NoError(t, err)
	pp.Println(info)
}
