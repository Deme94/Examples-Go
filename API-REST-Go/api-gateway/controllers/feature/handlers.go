package feature

import (
	"API-REST/api-gateway/controllers/feature/payloads"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/database/postgres/models/feature"
	psql "API-REST/services/database/postgres/predicates"
	"API-REST/services/logger"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *Controller) GetAll(ctx *fiber.Ctx) error {
	var features []*feature.Feature
	var err error
	// Query parameters
	predicates := psql.Predicates{}
	var queryParams payloads.QueryParams
	err = ctx.QueryParser(&queryParams)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}
	err = c.Validate.Struct(queryParams)
	if err != nil {
		return err
	}
	if queryParams.FromDate != nil && queryParams.ToDate != nil {
		predicates.Where("timestamp", ">=", queryParams.FromDate.Format("2006-01-02T15:04:05Z")).
			AndWhere("timestamp", "<=", queryParams.ToDate.Format("2006-01-02T15:04:05Z"))
		features, err = c.Model.GetAll(&predicates)
	} else {
		features, err = c.Model.GetAllMostRecent()
	}
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, features, "features")
}
func (c *Controller) GetByUserID(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Params("userId"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid user-id parameter"))
		return util.ErrorJSON(ctx, err)
	}

	// Query parameters
	predicates := psql.Predicates{}
	var queryParams payloads.QueryParams
	err = ctx.QueryParser(&queryParams)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}
	err = c.Validate.Struct(queryParams)
	if err != nil {
		return err
	}
	if queryParams.FromDate != nil && queryParams.ToDate != nil {
		predicates.Where("timestamp", ">=", queryParams.FromDate.Format("2006-01-02T15:04:05Z")).
			AndWhere("timestamp", "<=", queryParams.ToDate.Format("2006-01-02T15:04:05Z")).
			AndWhere("user_id", "=", userID.String())
		features, err := c.Model.GetAll(&predicates)
		if err != nil {
			return util.ErrorJSON(ctx, err)
		}
		return util.WriteJSON(ctx, http.StatusOK, features, "feature")
	} else {
		feature, err := c.Model.GetMostRecentByUserID(userID)
		if err != nil {
			return util.ErrorJSON(ctx, err)
		}
		featureResponse := payloads.GetMostRecentByUserIDResponse{Geom: feature.Geom, Timestamp: feature.Timestamp}
		return util.WriteJSON(ctx, http.StatusOK, featureResponse, "feature")
	}
}
func (c *Controller) Insert(ctx *fiber.Ctx) error {
	claimerID := uuid.MustParse(ctx.Locals("Claimer-ID").(string))

	var req payloads.InsertRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}
	err = c.Validate.Struct(req)
	if err != nil {
		return err
	}

	err = c.Model.Insert(&feature.Feature{Geom: req.Geom, Timestamp: req.Timestamp, UserID: claimerID})
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, "feature created successfully", "response")
}
