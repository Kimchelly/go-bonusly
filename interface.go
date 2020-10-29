package bonusly

import "context"

type Client interface {
	// CreateBonus creates a new bonus.
	CreateBonus(ctx context.Context, opts CreateBonusRequest) (*BonusResponse, error)
	// GetBonus gets a bonus by ID.
	GetBonus(ctx context.Context, id string) (*BonusResponse, error)
	// UpdateBonus updates a bonus by ID.
	UpdateBonus(ctx context.Context, id, reason string) (*BonusResponse, error)
	// UpdateBonus deletes a bonus by ID.
	DeleteBonus(ctx context.Context, id string) error
	// MyUserInfo returns information about the user making requests.
	MyUserInfo(ctx context.Context) (*UserInfoResponse, error)
	Close(ctx context.Context) error
}
