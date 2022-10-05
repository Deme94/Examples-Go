package controllers

import (
	util "API-REST/cmd/api/utilities"
	m "API-REST/models"
	"encoding/json"
	"log"
	"net/http"
	"time"

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
// ...

// PAYLOADS (json input) ----------------------------------------------------------------
type assetRequest struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

// ...

// API HANDLERS ---------------------------------------------------------------
func (c *AssetController) GetAll(w http.ResponseWriter, r *http.Request) {

}
func (c *AssetController) Get(w http.ResponseWriter, r *http.Request) {

}
func (c *AssetController) Insert(w http.ResponseWriter, r *http.Request) {
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

	err = c.model.Insert(&m.Asset{Name: req.Name, Date: date})
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, "asset inserted successfully", "response")
}
func (c *AssetController) Update(w http.ResponseWriter, r *http.Request) {

}
func (c *AssetController) Delete(w http.ResponseWriter, r *http.Request) {

}

// ...
