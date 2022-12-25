package asset

import (
	util "API-REST/api-gateway/utilities"
	"API-REST/services/database/models/asset"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CONTROLLER ***************************************************************
type Controller struct {
	Model *asset.Model
}

// METHODS CONTROLLER ---------------------------------------------------------------
func (c *Controller) getDateRangeFromQuery(query url.Values) (time.Time, time.Time, error) {
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

// PAYLOADS (json input and output) ----------------------------------------------------------------
type assetRequest struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type assetNameResponse struct {
	Name string `bson:"name"`
}

// ...

// API HANDLERS ---------------------------------------------------------------
func (c *Controller) GetAll(ctx *gin.Context) {
	fromDate, toDate, err := c.getDateRangeFromQuery(ctx.Request.URL.Query())
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	filterOptions := make(map[string]interface{})
	name := ctx.Query("name")
	if len(name) != 0 {
		filterOptions["name"] = name
	}
	// other filter options...

	assets, err := c.Model.GetAll(fromDate, toDate, filterOptions)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	util.WriteJSON(ctx, http.StatusOK, assets, "assets")
}
func (c *Controller) Get(ctx *gin.Context) {
	id := ctx.Param("id")

	a, err := c.Model.Get(id)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	util.WriteJSON(ctx, http.StatusOK, a, "asset")
}
func (c *Controller) GetWithAttributes(ctx *gin.Context) {
	id := ctx.Param("id")

	a, err := c.Model.GetWithAttributes(id)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	util.WriteJSON(ctx, http.StatusOK, a, "asset")
}
func (c *Controller) GetNames(ctx *gin.Context) {
	fromDate, toDate, err := c.getDateRangeFromQuery(ctx.Request.URL.Query())
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	assets, err := c.Model.GetNames(fromDate, toDate)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	var res []assetNameResponse
	var assetName assetNameResponse
	for _, a := range assets {
		assetName.Name = a.Name
		res = append(res, assetName)
	}

	util.WriteJSON(ctx, http.StatusOK, res, "assets")
}
func (c *Controller) Insert(ctx *gin.Context) {
	var req []assetRequest
	var assets []*asset.Asset

	err := ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	for _, a := range req {
		date, err := time.Parse("2006-01-02", a.Date)
		if err != nil {
			util.ErrorJSON(ctx, err)
			return
		}
		assets = append(assets, &asset.Asset{Name: a.Name, Date: date})
	}

	if len(assets) == 1 {
		err = c.Model.Insert(assets[0])
		if err != nil {
			util.ErrorJSON(ctx, err)
			return
		}
	} else {
		err = c.Model.InsertMany(assets)
		if err != nil {
			util.ErrorJSON(ctx, err)
			return
		}
	}

	util.WriteJSON(ctx, http.StatusOK, "assets inserted successfully", "response")
}
func (c *Controller) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req assetRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	var a asset.Asset
	a.ID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}
	a.Name = req.Name
	a.Date = date

	err = c.Model.Update(&a)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	util.WriteJSON(ctx, http.StatusOK, "asset updated successfully", "response")
}
func (c *Controller) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.Model.Delete(id)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	util.WriteJSON(ctx, http.StatusOK, "asset deleted successfully", "response")
}

// ...
