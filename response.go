package bonusly

import "time"

type CommonResponse struct {
	Success *bool   `json:"success,omitempty"`
	Message *string `json:"message,omitempty"`
}

type bonusResponseWrapper struct {
	CommonResponse
	Result BonusResponse `json:"result,omitempty"`
}

type bonusesResponseWrapper struct {
	CommonResponse
	Result []BonusResponse `json:"result,omitempty"`
}

type BonusResponse struct {
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
	ChildBonuses []BonusResponse `json:"child_bonuses,omitempty"`
	Via          *string         `json:"via,omitempty"`
	FamilyAmount *int            `json:"family_amount,omitempty"`
}

type rewardsResponseWrapper struct {
	CommonResponse
	Result []RewardsResponse `json:"result,omitempty"`
}

type RewardsResponse struct {
	Type    string           `json:"type,omitempty"`
	Name    string           `json:"name,omitempty"`
	Rewards []RewardResponse `json:"rewards,omitempty"`
}

type RewardResponse struct {
	Name                string                        `json:"name,omitempty"`
	ImageURL            string                        `json:"image_url,omitempty"`
	MinimumDisplayPrice string                        `json:"minimum_display_price,omitempty"`
	Description         RewardDescriptionResponse     `json:"description,omitempty"`
	DisclaimerHTML      string                        `json:"disclaimer_html,omitempty"`
	Categories          []string                      `json:"categories,omitempty"`
	Denominations       []RewardDenominationsResponse `json:"denominations,omitempty"`
}

type RewardDescriptionResponse struct {
	Text string `json:"text,omitempty"`
	HTML string `json:"html,omitempty"`
}

type RewardDenominationsResponse struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Price        int    `json:"price,omitempty"`
	DisplayPrice string `json:"display_price,omitempty"`
}

type userInfoResponseWrapper struct {
	CommonResponse
	Result UserInfoResponse `json:"result,omitempty"`
}

type UserInfoResponse struct {
	ID                           *string                `json:"id,omitempty"`
	UserName                     *string                `json:"username,omitempty"`
	Email                        *string                `json:"email,omitempty"`
	FirstName                    *string                `json:"first_name,omitempty"`
	LastName                     *string                `json:"last_name,omitempty"`
	ShortName                    *string                `json:"short_name,omitempty"`
	DisplayName                  *string                `json:"display_name,omitempty"`
	Path                         *string                `json:"path,omitempty"`
	FullPictureURL               *string                `json:"full_pic_url,omitempty"`
	LastActiveAt                 *time.Time             `json:"last_active_at,omitempty"`
	CreatedAt                    *time.Time             `json:"created_at,omitempty"`
	ExternalUniqueID             *string                `json:"external_unique_id,omitempty"`
	BudgetBoost                  *int                   `json:"budget_boost,omitempty"`
	UserMode                     *string                `json:"user_mode,omitempty"`
	TimeZone                     *string                `json:"time_zone,omitempty"`
	CanGive                      *bool                  `json:"can_give,omitempty"`
	CanReceive                   *bool                  `json:"can_receive,omitempty"`
	GiveAmounts                  *[]int                 `json:"give_amounts,omitempty"`
	Status                       *string                `json:"status,omitempty"`
	HiredOne                     *time.Time             `json:"hired_one,omitempty"`
	ManagerEmail                 *string                `json:"manager_email,omitempty"`
	EarningBalance               *int                   `json:"earning_balance,omitempty"`
	EarningBalanceWithCurrency   *string                `json:"earning_balance_with_currency,omitempty"`
	LifetimeEarnings             *int                   `json:"lifetime_earnings,omitempty"`
	LifetimeEarningsWithCurrency *string                `json:"lifetime_earnings_with_currency,omitempty"`
	GivingBalance                *int                   `json:"giving_balance,omitempty"`
	GivingBalanceWithCurrency    *string                `json:"giving_balance_with_currency,omitempty"`
	Department                   *string                `json:"department,omitempty"`
	Location                     *string                `json:"location,omitempty"`
	Number                       *string                `json:"number,omitempty"`
	CustomProperties             map[string]interface{} `json:"custom_properties,omitempty"`
}
