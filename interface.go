package bonusly

import "context"

type Client interface {
	// CreateBonus creates a new bonus.
	CreateBonus(ctx context.Context, req CreateBonusRequest) (*BonusResponse, error)
	// GetBonus gets a bonus by ID.
	GetBonus(ctx context.Context, id string) (*BonusResponse, error)
	// ListBonuses finds all bonuses matching the given request parameters.
	ListBonuses(ctx context.Context, req ListBonusesRequest) ([]BonusResponse, error)
	// UpdateBonus updates a bonus by ID.
	UpdateBonus(ctx context.Context, id, reason string) (*BonusResponse, error)
	// UpdateBonus deletes a bonus by ID.
	DeleteBonus(ctx context.Context, id string) error
	// ListRewards finds all rewards matching the given request parameters.
	ListRewards(ctx context.Context, req ListRewardsRequest) ([]RewardsResponse, error)
	// MyUserInfo returns information about the user making requests.
	MyUserInfo(ctx context.Context) (*UserInfoResponse, error)
	// Close closes the client and cleans up resources.
	Close(ctx context.Context) error
}
