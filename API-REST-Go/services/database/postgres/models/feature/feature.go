package feature

import (
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
)

type Feature struct {
	Geom      *postgis.PointS `json:"geom"`
	Timestamp time.Time       `json:"timestamp"`
	UserID    uuid.UUID       `json:"user_id"`
	// ...
}
