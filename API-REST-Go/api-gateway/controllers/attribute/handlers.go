package attribute

import (
	"API-REST/api-gateway/controllers/attribute/payloads"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/database/mongo/models/attribute"
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
	labelParam := ctx.Query("label")
	if len(nameParam) != 0 {
		filterOptions["name"] = nameParam
	}
	if len(labelParam) != 0 {
		filterOptions["label"] = labelParam
	}
	// other filter options...

	attributes, err := c.Model.GetAll(fromDate, toDate, filterOptions)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, attributes, "attributes")
}
func (c *Controller) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	a, err := c.Model.Get(id)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, a, "attribute")
}
func (c *Controller) Insert(ctx *fiber.Ctx) error {
	var req []payloads.AttributeRequest
	var attributes []*attribute.Attribute

	err := ctx.BodyParser(&req)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}
	err = c.Validate.Struct(req)
	if err != nil {
		return err
	}

	for _, a := range req {
		timestamp, err := time.Parse("2006-01-02T15:04:05", a.Timestamp)
		if err != nil {
			return util.ErrorJSON(ctx, err)
		}
		attributes = append(attributes,
			&attribute.Attribute{
				Metadata: attribute.AttributeMetadata{
					AssetName: a.AssetName,
					Name:      a.Name,
					Label:     a.Label,
					Unit:      a.Unit,
				},
				Timestamp: &timestamp,
				Value:     a.Value,
			},
		)
	}

	if len(attributes) == 1 {
		err = c.Model.Insert(attributes[0])
		if err != nil {
			return util.ErrorJSON(ctx, err)
		}
	} else {
		err = c.Model.InsertMany(attributes)
		if err != nil {
			return util.ErrorJSON(ctx, err)
		}
	}

	return util.WriteJSON(ctx, http.StatusOK, "attributes inserted successfully", "response")
}
func (c *Controller) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req payloads.AttributeRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}
	err = c.Validate.Struct(req)
	if err != nil {
		return err
	}

	timestamp, err := time.Parse("2006-01-02T15:04:05", req.Timestamp)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	var a attribute.Attribute
	a.ID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	a.Metadata = attribute.AttributeMetadata{
		AssetName: req.AssetName,
		Name:      req.Name,
		Label:     req.Label,
		Unit:      req.Unit,
	}
	a.Timestamp = &timestamp
	a.Value = req.Value

	err = c.Model.Update(&a)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, "attribute updated successfully", "response")
}
func (c *Controller) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.Model.Delete(id)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, "attribute deleted successfully", "response")
}
