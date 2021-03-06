package bonusly

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	accessToken := os.Getenv("BONUSLY_TOKEN")
	require.NotEmpty(t, accessToken, "missing access token")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpClient := getHTTPClient()
	defer putHTTPClient(httpClient)

	c, err := NewClient(ClientOptions{
		AccessToken: accessToken,
		HTTPClient:  httpClient,
	})
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, c.Close(ctx))
	}()

	t.Run("CreateBonus", func(t *testing.T) {
		t.Run("FailsWithBadInput", func(t *testing.T) {
			resp, err := c.CreateBonus(ctx, CreateBonusRequest{
				Reason: "+1 @nonexistent fail request",
			})
			assert.Error(t, err)
			assert.Zero(t, resp)
		})
	})
	t.Run("ListBonuses", func(t *testing.T) {
		t.Run("Succeeds", func(t *testing.T) {
			resp, err := c.ListBonuses(ctx, ListBonusesRequest{
				Limit: 1,
			})
			require.NoError(t, err)
			assert.NotZero(t, resp)
		})
		t.Run("FailsWithBadInput", func(t *testing.T) {
			resp, err := c.ListBonuses(ctx, ListBonusesRequest{
				GiverEmail: "nonexistent",
			})
			require.NoError(t, err)
			assert.NotZero(t, resp)
		})
	})
	t.Run("ListRewards", func(t *testing.T) {
		t.Run("Succeeds", func(t *testing.T) {
			resp, err := c.ListRewards(ctx, ListRewardsRequest{})
			require.NoError(t, err)
			assert.NotZero(t, resp)
		})
	})
	t.Run("MyUserInfo", func(t *testing.T) {
		t.Run("Succeeds", func(t *testing.T) {
			info, err := c.MyUserInfo(ctx)
			require.NoError(t, err)
			assert.NotZero(t, info)
		})
	})
}

func TestClientWithMockServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpClient := getHTTPClient()
	defer putHTTPClient(httpClient)

	c, err := NewClient(ClientOptions{
		AccessToken: "access_token",
		HTTPClient:  httpClient,
		BaseURL:     "https://private-anon-7c8306e7d6-bonusly.apiary-mock.com/api/v1",
	})
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, c.Close(ctx))
	}()

	t.Run("CreateBonus", func(t *testing.T) {
		t.Run("Succeeds", func(t *testing.T) {
			resp, err := c.CreateBonus(ctx, CreateBonusRequest{})
			require.NoError(t, err)
			assert.NotZero(t, resp)
		})
	})
	// For some reason, the mock server does not include the CommonResponse and
	// just lists the mock bonuses/reward, which differs from responses from the
	// production server.
	// t.Run("ListRewards", func(t *testing.T) {
	//     t.Run("Succeeds", func(t *testing.T) {
	//         resp, err := c.ListRewards(ctx, ListRewardsRequest{})
	//         require.NoError(t, err)
	//         assert.NotZero(t, resp)
	//     })
	// })
	// t.Run("ListBonuses", func(t *testing.T) {
	//     t.Run("Succeeds", func(t *testing.T) {
	//         resp, err := c.ListBonuses(ctx, ListBonusesRequest{
	//             Limit: 1,
	//         })
	//         require.NoError(t, err)
	//         assert.NotZero(t, resp)
	//     })
	// })
	t.Run("MyUserInfo", func(t *testing.T) {
		info, err := c.MyUserInfo(ctx)
		require.NoError(t, err)
		assert.NotZero(t, info)
	})
}
