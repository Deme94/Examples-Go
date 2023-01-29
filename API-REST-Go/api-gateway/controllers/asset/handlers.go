package asset

import (
	"API-REST/api-gateway/controllers/asset/payloads"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/database/mongo/models/asset"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *Controller) GetAll(ctx *fiber.Ctx) error {

	// Query parameters
	filterOptions := make(map[string]interface{})
	fromDate, toDate, err := c.getDateRange(ctx.Query("fromDate"), ctx.Query("toDate"))
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}
	nameParam := ctx.Query("name")
	if len(nameParam) != 0 {
		filterOptions["name"] = nameParam
	}
	// other filter options...

	assets, err := c.Model.GetAll(fromDate, toDate, filterOptions)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, assets, "assets")
}
func (c *Controller) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	a, err := c.Model.Get(id)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, a, "asset")
}
func (c *Controller) GetWithAttributes(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	a, err := c.Model.GetWithAttributes(id)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, a, "asset")
}
func (c *Controller) GetNames(ctx *fiber.Ctx) error {
	fromDate, toDate, err := c.getDateRange(ctx.Query("fromDate"), ctx.Query("toDate"))
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	assets, err := c.Model.GetNames(fromDate, toDate)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	var res []payloads.AssetNameResponse
	var assetName payloads.AssetNameResponse
	for _, a := range assets {
		assetName.Name = a.Name
		res = append(res, assetName)
	}

	return util.WriteJSON(ctx, http.StatusOK, res, "assets")
}
func (c *Controller) Insert(ctx *fiber.Ctx) error {
	var req []payloads.AssetRequest
	var assets []*asset.Asset

	err := ctx.BodyParser(&req)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	for _, a := range req {
		date, err := time.Parse("2006-01-02", a.Date)
		if err != nil {
			return util.ErrorJSON(ctx, err)
		}
		assets = append(assets, &asset.Asset{Name: a.Name, Date: &date})
	}

	if len(assets) == 1 {
		err = c.Model.Insert(assets[0])
		if err != nil {
			return util.ErrorJSON(ctx, err)
		}
	} else {
		err = c.Model.InsertMany(assets)
		if err != nil {
			return util.ErrorJSON(ctx, err)
		}
	}

	return util.WriteJSON(ctx, http.StatusOK, "assets inserted successfully", "response")
}
func (c *Controller) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req payloads.AssetRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	var a asset.Asset
	a.ID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}
	a.Name = req.Name
	a.Date = &date

	err = c.Model.Update(&a)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, "asset updated successfully", "response")
}
func (c *Controller) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.Model.Delete(id)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, "asset deleted successfully", "response")
}
