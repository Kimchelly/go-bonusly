package bonusly

import "context"

// MockClient is a mock Bonusly client that implements the Client interface.
type MockClient struct {
	CreateBonusRequest  CreateBonusRequest
	CreateBonusResponse BonusResponse

	GetBonusID       string
	GetBonusResponse BonusResponse

	UpdateBonusID       string
	UpdateBonusReason   string
	UpdateBonusResponse BonusResponse

	DeleteBonusID string

	MyUserInfoResponse UserInfoResponse
}

// CreateBonus records the CreateBonusRequest input and returns the mock
// client's CreateBonusResponse.
func (c *MockClient) CreateBonus(_ context.Context, req CreateBonusRequest) (*BonusResponse, error) {
	c.CreateBonusRequest = req
	return &c.CreateBonusResponse, nil
}

// CreateBonus records the bonus ID input and returns the mock client's
// GetBonusResponse.
func (c *MockClient) GetBonus(_ context.Context, id string) (*BonusResponse, error) {
	c.GetBonusID = id
	return &c.GetBonusResponse, nil
}

// CreateBonus records the bonus ID input and reason and returns the mock
// client's UpdateBonusResponse.
func (c *MockClient) UpdateBonus(_ context.Context, id, reason string) (*BonusResponse, error) {
	c.UpdateBonusID = id
	c.UpdateBonusReason = reason
	return &c.UpdateBonusResponse, nil
}

// DeleteBonus records the bonus ID input.
func (c *MockClient) DeleteBonus(_ context.Context, id string) error {
	c.DeleteBonusID = id
	return nil
}

// MyUserInfo returns the mock client's MyUserInfoResponse.
func (c *MockClient) MyUserInfo(_ context.Context) (*UserInfoResponse, error) {
	return &c.MyUserInfoResponse, nil
}

// Close is a no-op.
func (c *MockClient) Close(_ context.Context) error {
	return nil
}
