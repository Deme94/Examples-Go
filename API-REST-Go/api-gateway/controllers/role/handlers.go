package role

import (
	"API-REST/api-gateway/controllers/role/payloads"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/database/postgres/models/role"
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

	roles, err := c.Model.GetAll(&predicates)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	var all []*payloads.GetResponse
	for _, role := range roles {
		all = append(all, &payloads.GetResponse{
			ID:   role.ID,
			Name: role.Name,
		})
	}

	allResponse := payloads.GetAllResponse{Roles: all}
	return util.WriteJSON(ctx, http.StatusOK, allResponse)
}
func (c *Controller) Get(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		return util.ErrorJSON(ctx, err)
	}

	role, err := c.Model.Get(id)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, role, "role")
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

	err = c.Model.Insert(&role.Role{Name: req.Name})
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, "role created successfully", "response")
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

	var role role.Role
	role.ID = id
	role.Name = req.Name

	err = c.Model.Update(&role)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, true, "OK")
}
func (c *Controller) UpdatePermissions(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		return util.ErrorJSON(ctx, err)
	}

	var req payloads.UpdatePermissionsRequest

	err = ctx.BodyParser(&req)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}
	err = c.Validate.Struct(req)
	if err != nil {
		return err
	}

	err = c.Model.UpdatePermissions(id, req.PermissionIDs...)
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
