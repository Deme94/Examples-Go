package payloads

type AssetRequest struct {
	Name string `json:"name,omitempty"`
	Date string `json:"date,omitempty"`
}
