package attribute

import (
	"API-REST/api-gateway/controllers/attribute/payloads"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/database/mongo/models/attribute"
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
		util.ErrorJSON(ctx, err)
		return
	}

	util.WriteJSON(ctx, http.StatusOK, attributes, "attributes")
}
func (c *Controller) Get(ctx *gin.Context) {
	id := ctx.Param("id")

	a, err := c.Model.Get(id)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	util.WriteJSON(ctx, http.StatusOK, a, "attribute")
}
func (c *Controller) Insert(ctx *gin.Context) {
	var req []payloads.AttributeRequest
	var attributes []*attribute.Attribute

	err := ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	for _, a := range req {
		timestamp, err := time.Parse("2006-01-02T15:04:05", a.Timestamp)
		if err != nil {
			util.ErrorJSON(ctx, err)
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
				Timestamp: &timestamp,
				Value:     a.Value,
			},
		)
	}

	if len(attributes) == 1 {
		err = c.Model.Insert(attributes[0])
		if err != nil {
			util.ErrorJSON(ctx, err)
			return
		}
	} else {
		err = c.Model.InsertMany(attributes)
		if err != nil {
			util.ErrorJSON(ctx, err)
			return
		}
	}

	util.WriteJSON(ctx, http.StatusOK, "attributes inserted successfully", "response")
}
func (c *Controller) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req payloads.AttributeRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	timestamp, err := time.Parse("2006-01-02T15:04:05", req.Timestamp)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	var a attribute.Attribute
	a.ID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
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
		util.ErrorJSON(ctx, err)
		return
	}

	util.WriteJSON(ctx, http.StatusOK, "attribute updated successfully", "response")
}
func (c *Controller) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.Model.Delete(id)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	util.WriteJSON(ctx, http.StatusOK, "attribute deleted successfully", "response")
}
