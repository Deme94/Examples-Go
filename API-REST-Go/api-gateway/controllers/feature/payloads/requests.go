package payloads

import (
	"time"

	"github.com/cridenour/go-postgis"
)

type QueryParams struct {
	FromDate *time.Time `query:"fromDate" validate:"required_with=ToDate"`
	ToDate   *time.Time `query:"toDate" validate:"required_with=FromDate"`
}

type InsertRequest struct {
	Geom      *postgis.PointS `json:"geom" validate:"required"`
	Timestamp time.Time       `json:"timestamp" validate:"required"`
}
