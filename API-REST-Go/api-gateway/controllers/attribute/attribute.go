package attribute

import (
	util "API-REST/api-gateway/utilities"
	"API-REST/services/database/models/attribute"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CONTROLLER ***************************************************************
type Controller struct {
	Model *attribute.Model
}

// METHODS CONTROLLER ---------------------------------------------------------------
func (c *Controller) getDateRangeFromQuery(query url.Values) (time.Time, time.Time, error) {
	var fromDate time.Time
	var toDate time.Time

	// if query parameters
	fromDateString := query.Get("fromDate")
	toDateString := query.Get("toDate")
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

// PAYLOADS (json input and output) ----------------------------------------------------------------
type attributeRequest struct {
	AssetName string  `json:"asset_name"`
	Name      string  `json:"name"`
	Label     string  `json:"label"`
	Unit      string  `json:"unit"`
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
}

// ...

// API HANDLERS ---------------------------------------------------------------
func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	fromDate, toDate, err := c.getDateRangeFromQuery(r.URL.Query())
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	filterOptions := make(map[string]interface{})
	name := r.URL.Query().Get("name")
	label := r.URL.Query().Get("label")
	if len(name) != 0 {
		filterOptions["name"] = name
	}
	if len(label) != 0 {
		filterOptions["label"] = label
	}
	// other filter options...

	attributes, err := c.Model.GetAll(fromDate, toDate, filterOptions)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, attributes, "attributes")
}
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id := params.ByName("id")

	a, err := c.Model.Get(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, a, "attribute")
}
func (c *Controller) Insert(w http.ResponseWriter, r *http.Request) {
	var req []attributeRequest
	var attributes []*attribute.Attribute

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	for _, a := range req {
		timestamp, err := time.Parse("2006-01-02T15:04:05", a.Timestamp)
		if err != nil {
			util.ErrorJSON(w, err)
			return
		}
		attributes = append(attributes,
			&attribute.Attribute{
				Metadata: attribute.AttributeMetadata{
					AssetName: a.AssetName,
					Name:      a.Name,
					Label:     a.Label,
					Unit:      a.Unit,
				},
				Timestamp: timestamp,
				Value:     a.Value,
			},
		)
	}

	if len(attributes) == 1 {
		err = c.Model.Insert(attributes[0])
		if err != nil {
			util.ErrorJSON(w, err)
			return
		}
	} else {
		err = c.Model.InsertMany(attributes)
		if err != nil {
			util.ErrorJSON(w, err)
			return
		}
	}

	util.WriteJSON(w, http.StatusOK, "attributes inserted successfully", "response")
}
func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	var req attributeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	timestamp, err := time.Parse("2006-01-02T15:04:05", req.Timestamp)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	var a attribute.Attribute
	a.ID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	a.Metadata = attribute.AttributeMetadata{
		AssetName: req.AssetName,
		Name:      req.Name,
		Label:     req.Label,
		Unit:      req.Unit,
	}
	a.Timestamp = timestamp
	a.Value = req.Value

	err = c.Model.Update(&a)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, "attribute updated successfully", "response")
}
func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	err := c.Model.Delete(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, "attribute deleted successfully", "response")
}

// ...
