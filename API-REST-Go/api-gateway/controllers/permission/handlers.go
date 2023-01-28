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

	"github.com/gin-gonic/gin"
)

func (c *Controller) GetAll(ctx *gin.Context) {

	// Query parameters
	predicates := psql.Predicates{}
	pageParam := ctx.Query("page")
	pageSizeParam := ctx.Query("pageSize")
	nameParam := ctx.Query("name")
	if len(pageParam) != 0 && len(pageSizeParam) != 0 {
		page, err := strconv.Atoi(pageParam)
		if err != nil {
			util.ErrorJSON(ctx, err)
			return
		}
		pageSize, err := strconv.Atoi(pageSizeParam)
		if err != nil {
			util.ErrorJSON(ctx, err)
			return
		}
		if page == 0 || pageSize == 0 {
			util.ErrorJSON(ctx, errors.New("page and pageSize params must be greater than 0"))
			return
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
		util.ErrorJSON(ctx, err)
		return
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
	util.WriteJSON(ctx, http.StatusOK, allResponse, "response")
}
func (c *Controller) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		util.ErrorJSON(ctx, err)
		return
	}

	permission, err := c.Model.Get(id)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	util.WriteJSON(ctx, http.StatusOK, permission, "permission")
}
func (c *Controller) Insert(ctx *gin.Context) {
	var req payloads.InsertRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	err = c.Model.Insert(&permission.Permission{
		Resource:  req.Resource,
		Operation: req.Operation,
	})
	if err != nil {
		util.ErrorJSON(ctx, err)
		return
	}

	util.WriteJSON(ctx, http.StatusOK, "permission created successfully", "response")
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

	var permission permission.Permission
	permission.ID = id
	permission.Resource = req.Resource
	permission.Operation = req.Operation

	err = c.Model.Update(&permission)
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