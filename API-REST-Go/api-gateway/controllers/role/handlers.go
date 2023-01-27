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

	"github.com/gin-gonic/gin"
)

func (c *Controller) GetAll(ctx *gin.Context) {

	// Query parameters
	nameParam := ctx.Query("name")
	predicates := psql.Predicates{}
	if len(nameParam) != 0 {
		predicates.Where("name", "=", nameParam)
	}

	roles, err := c.Model.GetAll(&predicates)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	var all []*payloads.GetResponse
	for _, role := range roles {
		all = append(all, &payloads.GetResponse{
			ID:   role.ID,
			Name: role.Name,
		})
	}

	allResponse := payloads.GetAllResponse{Roles: all}
	util.WriteJSON(ctx, http.StatusOK, allResponse, "response")
}
func (c *Controller) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		util.ErrorJSON(ctx, err)
		return
	}

	role, err := c.Model.Get(id)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	util.WriteJSON(ctx, http.StatusOK, role, "role")
}
func (c *Controller) Insert(ctx *gin.Context) {
	var req payloads.InsertRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	err = c.Model.Insert(&role.Role{Name: req.Name})
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	util.WriteJSON(ctx, http.StatusOK, "role created successfully", "response")
}
func (c *Controller) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		util.ErrorJSON(ctx, err)
		return
	}

	var req payloads.UpdateRequest

	err = ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	var role role.Role
	role.ID = id
	role.Name = req.Name

	err = c.Model.Update(&role)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	ok := payloads.OkResponse{
		OK: true,
	}
	util.WriteJSON(ctx, http.StatusOK, ok, "OK")
}
func (c *Controller) UpdatePermissions(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		util.ErrorJSON(ctx, err)
		return
	}

	var req payloads.UpdatePermissionsRequest

	err = ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	err = c.Model.UpdatePermissions(id, req.PermissionIDs...)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	ok := payloads.OkResponse{
		OK: true,
	}
	util.WriteJSON(ctx, http.StatusOK, ok, "OK")
}
func (c *Controller) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		util.ErrorJSON(ctx, err)
		return
	}

	err = c.Model.Delete(id)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	ok := payloads.OkResponse{
		OK: true,
	}
	util.WriteJSON(ctx, http.StatusOK, ok, "OK")
}
