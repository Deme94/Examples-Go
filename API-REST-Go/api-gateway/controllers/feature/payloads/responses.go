package payloads

import (
	"time"

	"github.com/cridenour/go-postgis"
)

// Default responses from models is better in this case for performance reasons

type GetMostRecentByUserIDResponse struct {
	Geom      *postgis.PointS `json:"geom"`
	Timestamp time.Time       `json:"timestamp"`
}
