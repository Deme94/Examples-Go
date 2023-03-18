package payloads

type AttributeRequest struct {
	AssetName string  `json:"asset_name,omitempty"`
	Name      string  `json:"name,omitempty"`
	Label     string  `json:"label,omitempty"`
	Unit      string  `json:"unit,omitempty"`
	Timestamp string  `json:"timestamp,omitempty"`
	Value     float64 `json:"value,omitempty"`
}
