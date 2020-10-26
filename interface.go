package bonusly

import "context"

type Client interface {
	// CreateBonus creates a new bonus.
	CreateBonus(ctx context.Context, opts CreateBonusRequest) (*CreateBonusResponse, error)
	// MyUserInfo returns information about the user making requests.
	MyUserInfo(ctx context.Context) (*UserInfoResponse, error)
	Close(ctx context.Context) error
}
