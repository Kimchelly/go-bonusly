package bonusly

import "context"

// MockClient is a mock Bonusly client that implements the Client interface.
type MockClient struct {
	CreateBonusInput  CreateBonusRequest
	CreateBonusOutput CreateBonusResponse

	MyUserInfoOutput UserInfoResponse
}

// CreateBonus records the CreateBonusRequest input and returns the mock
// client's CreateBonusOutput.
func (c *MockClient) CreateBonus(_ context.Context, in CreateBonusRequest) (*CreateBonusResponse, error) {
	c.CreateBonusInput = in
	return &c.CreateBonusOutput, nil
}

// CreateBonus returns the mock client's MyUserInfoOutput.
func (c *MockClient) MyUserInfo(_ context.Context) (*UserInfoResponse, error) {
	return &c.MyUserInfoOutput, nil
}

// Close is a no-op.
func (c *MockClient) Close(_ context.Context) error {
	return nil
}
