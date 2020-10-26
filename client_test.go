package bonusly

import (
	"context"
	"os"
	"testing"

	"github.com/evergreen-ci/utility"
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

	t.Run("MyUserInfo", func(t *testing.T) {
		info, err := c.MyUserInfo(ctx)
		require.NoError(t, err)
		assert.NotZero(t, info)
	})

	t.Run("CreateBonus", func(t *testing.T) {
		t.Run("FailsWithBadUsername", func(t *testing.T) {
			resp, err := c.CreateBonus(ctx, CreateBonusRequest{
				Reason: "+1 @nonexistent fail request",
			})
			assert.Error(t, err)
			assert.Zero(t, resp)
		})
	})
}
