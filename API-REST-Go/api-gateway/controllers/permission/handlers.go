package permission

import (
	"API-REST/api-gateway/controllers/permission/payloads"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/database/postgres/models/permission"
	psql "API-REST/services/database/postgres/predicates"
	"API-REST/services/logger"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (c *Controller) GetAll(ctx *fiber.Ctx) error {

	// Query parameters
	predicates := psql.Predicates{}
	pageParam := ctx.Query("page")
	pageSizeParam := ctx.Query("pageSize")
	nameParam := ctx.Query("name")
	if len(pageParam) != 0 && len(pageSizeParam) != 0 {
		page, err := strconv.Atoi(pageParam)
		if err != nil {
			return util.ErrorJSON(ctx, err)
		}
		pageSize, err := strconv.Atoi(pageSizeParam)
		if err != nil {
			return util.ErrorJSON(ctx, err)
		}
		if page == 0 || pageSize == 0 {
			return util.ErrorJSON(ctx, errors.New("page and pageSize params must be greater than 0"))
		}
		limit := pageSize
		offset := (page - 1) * pageSize
		predicates.Offset(offset).Limit(limit)
	}
	if len(nameParam) != 0 {
		predicates.Where("name", "=", nameParam)
	}

	permissions, err := c.Model.GetAll(&predicates)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	var all []*payloads.GetResponse
	for _, permission := range permissions {
		all = append(all, &payloads.GetResponse{
			ID:        permission.ID,
			Resource:  permission.Resource,
			Operation: permission.Operation,
		})
	}

	allResponse := payloads.GetAllResponse{Permissions: all}
	return util.WriteJSON(ctx, http.StatusOK, allResponse)
}
func (c *Controller) Get(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		return util.ErrorJSON(ctx, err)
	}

	permission, err := c.Model.Get(id)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, permission, "permission")
}
func (c *Controller) Insert(ctx *fiber.Ctx) error {
	var req payloads.InsertRequest

	err := ctx.BodyParser(&req)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}
	err = c.Validate.Struct(req)
	if err != nil {
		return err
	}

	err = c.Model.Insert(&permission.Permission{
		Resource:  req.Resource,
		Operation: req.Operation,
	})
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, "permission created successfully", "response")
}
func (c *Controller) Update(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		return util.ErrorJSON(ctx, err)
	}

	var req payloads.UpdateRequest

	err = ctx.BodyParser(&req)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}
	err = c.Validate.Struct(req)
	if err != nil {
		return err
	}

	var permission permission.Permission
	permission.ID = id
	permission.Resource = req.Resource
	permission.Operation = req.Operation

	err = c.Model.Update(&permission)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, true, "OK")
}
func (c *Controller) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		return util.ErrorJSON(ctx, err)
	}

	err = c.Model.Delete(id)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, true, "OK")
}
