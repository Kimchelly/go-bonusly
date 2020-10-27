package bonusly

const (
	productionBaseURL = "https://bonus.ly/api/v1"

	contentType = "application/json"
)

type CreateBonusRequest struct {
	GiverEmail    string `json:"giver_email,omitempty"`
	Reason        string `json:"reason,omitempty"`
	ParentBonusID string `json:"parent_bonus_id,omitempty"`
}

// TODO: figure out which parameters are required
func (r *CreateBonusRequest) Validate() error {
	catcher := newBasicCatcher()
	return catcher.Resolve()
}
