package payloads

type AttributeRequest struct {
	AssetName string  `json:"asset_name"`
	Name      string  `json:"name"`
	Label     string  `json:"label"`
	Unit      string  `json:"unit"`
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
}
