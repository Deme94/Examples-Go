package asset

import (
	"API-REST/services/database/mongo/models/asset"
	"net/url"
	"time"
)

type Controller struct {
	Model *asset.Model
}

// METHODS CONTROLLER ---------------------------------------------------------------
func (c *Controller) getDateRangeFromQuery(query url.Values) (time.Time, time.Time, error) {
	var fromDate time.Time
	var toDate time.Time

	// Query parameters
	fromDateString := query.Get("fromDate")
	toDateString := query.Get("toDate")
	if len(fromDateString) != 0 {
		from, err := time.Parse("2006-01-02", fromDateString)
		if err != nil {
			return fromDate, toDate, err
		}
		fromDate = from
	}
	if len(toDateString) != 0 {
		to, err := time.Parse("2006-01-02", toDateString)
		if err != nil {
			return fromDate, toDate, err
		}
		toDate = to
	}

	return fromDate, toDate, nil
}

// ...
