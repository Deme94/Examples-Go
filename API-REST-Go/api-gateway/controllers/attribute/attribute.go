package attribute

import (
	"API-REST/services/database/mongo/models/attribute"
	"time"
)

type Controller struct {
	Model *attribute.Model
}

// METHODS CONTROLLER ---------------------------------------------------------------
func (c *Controller) getDateRange(fromDateString string, toDateString string) (time.Time, time.Time, error) {
	var fromDate time.Time
	var toDate time.Time

	// Query parameters
	if len(fromDateString) != 0 {
		from, err := time.Parse("2006-01-02T15:04:05", fromDateString)
		if err != nil {
			return fromDate, toDate, err
		}
		fromDate = from
	}
	if len(toDateString) != 0 {
		to, err := time.Parse("2006-01-02T15:04:05", toDateString)
		if err != nil {
			return fromDate, toDate, err
		}
		toDate = to
	}

	return fromDate, toDate, nil
}

// ...
