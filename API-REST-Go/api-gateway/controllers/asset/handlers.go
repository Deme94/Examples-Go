package asset

import (
	"API-REST/api-gateway/controllers/asset/payloads"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/database/mongo/models/asset"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *Controller) GetAll(ctx *gin.Context) {

	// Query parameters
	filterOptions := make(map[string]interface{})
	fromDate, toDate, err := c.getDateRangeFromQuery(ctx.Request.URL.Query())
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}
	nameParam := ctx.Query("name")
	if len(nameParam) != 0 {
		filterOptions["name"] = nameParam
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

	var res []payloads.AssetNameResponse
	var assetName payloads.AssetNameResponse
	for _, a := range assets {
		assetName.Name = a.Name
		res = append(res, assetName)
	}

	util.WriteJSON(ctx, http.StatusOK, res, "assets")
}
func (c *Controller) Insert(ctx *gin.Context) {
	var req []payloads.AssetRequest
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
		assets = append(assets, &asset.Asset{Name: a.Name, Date: &date})
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

	var req payloads.AssetRequest
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
	a.Date = &date

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
