package bonusly

import "time"

type CommonResponse struct {
	Success *bool   `json:"success,omitempty"`
	Message *string `json:"message,omitempty"`
}

type createBonusResponseWrapper struct {
	CommonResponse
	Result CreateBonusResponse `json:"result,omitempty"`
}

type CreateBonusResponse struct {
	ID                 *string           `json:"id,omitempty"`
	CreatedAt          *time.Time        `json:"created_at,omitempty"`
	Reason             *string           `json:"reason,omitempty"`
	ReasonHTML         *string           `json:"reason_html,omitempty"`
	Amount             *int              `json:"amount,omitempty"`
	AmountWithCurrency *string           `json:"amount_with_currency,omitempty"`
	Value              string            `json:"value,omitempty"`
	Giver              *UserInfoResponse `json:"giver,omitempty"`
	Receiver           *UserInfoResponse `json:"receiver,omitempty"`
	ChildCount         *int              `json:"child_count,omitempty"`
	// TODO: figure out what the expected structure of this is.
	ChildBonuses []interface{} `json:"child_bonuses,omitempty"`
	Via          *string       `json:"via,omitempty"`
	FamilyAmount *int          `json:"family_amount,omitempty"`
}

type userInfoResponseWrapper struct {
	CommonResponse
	Result UserInfoResponse `json:"result,omitempty"`
}

type UserInfoResponse struct {
	ID             *string    `json:"id,omitempty"`
	UserName       *string    `json:"username,omitempty"`
	Email          *string    `json:"email,omitempty"`
	FirstName      *string    `json:"first_name,omitempty"`
	LastName       *string    `json:"last_name,omitempty"`
	ShortName      *string    `json:"short_name,omitempty"`
	DisplayName    *string    `json:"display_name,omitempty"`
	Path           *string    `json:"path,omitempty"`
	FullPictureURL *string    `json:"full_pic_url,omitempty"`
	LastActiveAt   *time.Time `json:"last_active_at,omitempty"`
}
