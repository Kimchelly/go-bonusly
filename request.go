package bonusly

import (
	"strconv"
	"time"
)

const (
	productionBaseURL = "https://bonus.ly/api/v1"

	contentType = "application/json"
)

type CreateBonusRequest struct {
	GiverEmail    string `json:"giver_email,omitempty"`
	Reason        string `json:"reason,omitempty"`
	ParentBonusID string `json:"parent_bonus_id,omitempty"`
}

type ListBonusesRequest struct {
	Limit              uint
	Skip               uint
	StartTime          time.Time
	EndTime            time.Time
	GiverEmail         string
	ReceiverEmail      string
	UserEmail          string
	HashTag            string
	IncludeChildren    bool
	CustomPropertyName string
	ShowPrivateBonuses bool
}

func (r *ListBonusesRequest) QueryMap() map[string]string {
	q := map[string]string{}
	if r.Limit != 0 {
		q["limit"] = strconv.Itoa(int(r.Limit))
	}
	if r.Skip != 0 {
		q["skip"] = strconv.Itoa(int(r.Skip))
	}
	if !r.StartTime.IsZero() {
		q["start_time"] = r.StartTime.UTC().String()
	}
	if !r.EndTime.IsZero() {
		q["end_time"] = r.EndTime.UTC().String()
	}
	if r.GiverEmail != "" {
		q["giver_email"] = r.GiverEmail
	}
	if r.ReceiverEmail != "" {
		q["receiver_email"] = r.ReceiverEmail
	}
	if r.UserEmail != "" {
		q["user_email"] = r.UserEmail
	}
	if r.HashTag != "" {
		q["hashtag"] = r.HashTag
	}
	if r.IncludeChildren {
		q["include_children"] = strconv.FormatBool(r.IncludeChildren)
	}
	if r.CustomPropertyName != "" {
		q["custom_property_name"] = r.CustomPropertyName
	}
	if r.ShowPrivateBonuses {
		q["show_private_bonuses"] = strconv.FormatBool(r.ShowPrivateBonuses)
	}
	return q
}

type ListRewardsRequest struct {
	CatalogCountry string
	RequestCountry string
	PersonalizeFor string
}

func (r *ListRewardsRequest) QueryMap() map[string]string {
	q := map[string]string{}
	if r.CatalogCountry != "" {
		q["catalog_country"] = r.CatalogCountry
	}
	if r.RequestCountry != "" {
		q["request_country"] = r.RequestCountry
	}
	if r.PersonalizeFor != "" {
		q["personalize_for"] = r.PersonalizeFor
	}
	return q
}
