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
type AssetController struct {
	model  *m.AssetModel
	logger *log.Logger
}

func NewAssetController(coll *mongo.Collection, logger *log.Logger) *AssetController {
	c := AssetController{}
	c.model = &m.AssetModel{Coll: coll}
	c.logger = logger

	return &c
}

// METHODS CONTROLLER ---------------------------------------------------------------
func (c *AssetController) getDateRangeFromQuery(query url.Values) (time.Time, time.Time, error) {
	var fromDate time.Time
	var toDate time.Time

	// if query parameters
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

// PAYLOADS (json input) ----------------------------------------------------------------
type assetRequest struct {
	Name string `bson:"name"`
	Date string `bson:"date"`
}

// ...

// API HANDLERS ---------------------------------------------------------------
func (c *AssetController) GetAll(w http.ResponseWriter, r *http.Request) {
	fromDate, toDate, err := c.getDateRangeFromQuery(r.URL.Query())
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	filterOptions := make(map[string]interface{})
	name := r.URL.Query().Get("name")
	if len(name) != 0 {
		filterOptions["name"] = name
	}
	// other filter options...

	assets, err := c.model.GetAll(fromDate, toDate, filterOptions)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, assets, "assets")
}
func (c *AssetController) Get(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id := params.ByName("id")

	a, err := c.model.Get(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, a, "asset")
}
func (c *AssetController) GetNames(w http.ResponseWriter, r *http.Request) {
	fromDate, toDate, err := c.getDateRangeFromQuery(r.URL.Query())
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	assets, err := c.model.GetNames(fromDate, toDate)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, assets, "assets")
}
func (c *AssetController) Insert(w http.ResponseWriter, r *http.Request) {
	var req []assetRequest
	var assets []*m.Asset

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	for _, a := range req {
		date, err := time.Parse("2006-01-02", a.Date)
		if err != nil {
			util.ErrorJSON(w, err)
			return
		}
		assets = append(assets, &m.Asset{Name: a.Name, Date: date})
	}

	if len(assets) == 1 {
		err = c.model.Insert(assets[0])
		if err != nil {
			util.ErrorJSON(w, err)
			return
		}
	} else {
		err = c.model.InsertMany(assets)
		if err != nil {
			util.ErrorJSON(w, err)
			return
		}
	}

	util.WriteJSON(w, http.StatusOK, "assets inserted successfully", "response")
}
func (c *AssetController) Update(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	var req assetRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	var a m.Asset
	a.ID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}
	a.Name = req.Name
	a.Date = date

	err = c.model.Update(&a)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, "asset updated successfully", "response")
}
func (c *AssetController) Delete(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	err := c.model.Delete(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, "asset deleted successfully", "response")
}

// ...
