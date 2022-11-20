package controllers

import (
	util "API-REST/cmd/api/utilities"
	m "API-REST/models"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CONTROLLER ***************************************************************
type AttributeController struct {
	model  *m.AttributeModel
	logger *log.Logger
}

func NewAttributeController(coll *mongo.Collection, db *mongo.Database, logger *log.Logger) *AttributeController {
	c := AttributeController{}
	c.model = &m.AttributeModel{Coll: coll, Db: db}
	c.logger = logger

	return &c
}

// METHODS CONTROLLER ---------------------------------------------------------------
func (c *AttributeController) getDateRangeFromQuery(query url.Values) (time.Time, time.Time, error) {
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
func (c *AttributeController) GetAll(w http.ResponseWriter, r *http.Request) {
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

	attributes, err := c.model.GetAll(fromDate, toDate, filterOptions)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, attributes, "attributes")
}
func (c *AttributeController) Get(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id := params.ByName("id")

	a, err := c.model.Get(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, a, "attribute")
}
func (c *AttributeController) Insert(w http.ResponseWriter, r *http.Request) {
	var req []attributeRequest
	var attributes []*m.Attribute

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
			&m.Attribute{
				Metadata: m.AttributeMetadata{
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
		err = c.model.Insert(attributes[0])
		if err != nil {
			util.ErrorJSON(w, err)
			return
		}
	} else {
		err = c.model.InsertMany(attributes)
		if err != nil {
			util.ErrorJSON(w, err)
			return
		}
	}

	util.WriteJSON(w, http.StatusOK, "attributes inserted successfully", "response")
}
func (c *AttributeController) Update(w http.ResponseWriter, r *http.Request) {
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

	var a m.Attribute
	a.ID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	a.Metadata = m.AttributeMetadata{
		AssetName: req.AssetName,
		Name:      req.Name,
		Label:     req.Label,
		Unit:      req.Unit,
	}
	a.Timestamp = timestamp
	a.Value = req.Value

	err = c.model.Update(&a)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, "attribute updated successfully", "response")
}
func (c *AttributeController) Delete(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	err := c.model.Delete(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, "attribute deleted successfully", "response")
}

// ...
