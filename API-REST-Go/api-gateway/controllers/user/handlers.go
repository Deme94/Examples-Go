package user

import (
	"API-REST/api-gateway/controllers/user/payloads"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/database/postgres/models/user"
	psql "API-REST/services/database/postgres/predicates"
	"API-REST/services/logger"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

func (c *Controller) GetAll(ctx *fiber.Ctx) error {

	// Query parameters
	predicates := psql.Predicates{}
	pageParam := ctx.Query("page")
	pageSizeParam := ctx.Query("pageSize")
	yearParam := ctx.Query("year")
	deletedParam := ctx.Query("deleted")
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
	if len(yearParam) != 0 {
		startDate := fmt.Sprint(yearParam, "-01-01")
		endDate := fmt.Sprint(yearParam, "-12-31")
		predicates.Where("created_at", ">=", startDate).AndWhere("created_at", "<=", endDate)
	}
	if len(deletedParam) != 0 {
		deleted, err := strconv.ParseBool(deletedParam)
		if err == nil {
			if predicates.HasWhere() {
				if deleted {
					predicates.AndWhereNotNull("deleted_at")
				} else {
					predicates.AndWhereNull("deleted_at")
				}
			} else {
				if deleted {
					predicates.WhereNotNull("deleted_at")
				} else {
					predicates.WhereNull("deleted_at")
				}
			}
		}
	}

	usrs, err := c.Model.GetAll(&predicates)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	var all []*payloads.GetResponse
	for _, user := range usrs {
		all = append(all, &payloads.GetResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		})
	}

	allResponse := payloads.GetAllResponse{Users: all}
	return util.WriteJSON(ctx, http.StatusOK, allResponse, "response")
}
func (c *Controller) Get(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		return util.ErrorJSON(ctx, err)
	}

	u, err := c.Model.Get(id)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, u, "user")
}
func (c *Controller) GetPhoto(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		return util.ErrorJSON(ctx, err)
	}

	imageName, err := c.Model.GetPhoto(id)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	filePath := "./storage/users/" + imageName + ".webp"
	return ctx.SendFile(filePath)
}
func (c *Controller) GetCV(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		return util.ErrorJSON(ctx, err)
	}

	cvName, err := c.Model.GetCV(id)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	filePath := "./storage/users/" + cvName + ".pdf"
	return ctx.SendFile(filePath)
}
func (c *Controller) Insert(ctx *fiber.Ctx) error {
	var req payloads.InsertRequest

	err := ctx.BodyParser(&req)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	hashedPassword, err := c.HashPassword(req.Password)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	var u user.User
	u.Username = req.Username
	u.Email = req.Email
	u.Password = hashedPassword
	u.Nick = req.Nick
	u.FirstName = req.FirstName
	u.LastName = req.LastName
	u.Phone = req.Phone
	u.Address = req.Address

	err = c.Model.Insert(&u)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, "user created successfully", "response")
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

	var u user.User
	u.ID = id
	u.Nick = req.Nick
	u.FirstName = req.FirstName
	u.LastName = req.LastName
	u.Phone = req.Phone
	u.Address = req.Address

	err = c.Model.Update(&u)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, true, "OK")
}
func (c *Controller) UpdateRoles(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		return util.ErrorJSON(ctx, err)
	}

	var req payloads.UpdateRolesRequest

	err = ctx.BodyParser(&req)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	err = c.Model.UpdateRoles(id, req.RoleIDs...)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, true, "OK")
}
func (c *Controller) UpdatePhoto(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		return util.ErrorJSON(ctx, err)
	}

	var req payloads.UpdatePhotoRequest

	err = ctx.BodyParser(&req)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	// Decode base64 webp to bytes
	unbased, err := base64.StdEncoding.DecodeString(req.PhotoBase64)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}
	// Decode bytes to image
	reader := bytes.NewReader(unbased)
	image, err := webp.Decode(reader, &decoder.Options{})
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}
	// Create our own file
	imageName := "user" + fmt.Sprint(id)
	f, err := os.OpenFile("./storage/users/"+imageName+".webp", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}
	defer f.Close()
	// Encode image into file
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	if err != nil {
		log.Fatalln(err)
	}
	webp.Encode(f, image, options)
	logger.Logger.Println("user's photo saved")

	err = c.Model.UpdatePhoto(id, imageName)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, true, "OK")
}
func (c *Controller) UpdateCV(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		return util.ErrorJSON(ctx, err)
	}

	// Retrieve the file
	file, err := ctx.FormFile("file")
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	// Check if file is .pdf
	if filepath.Ext(file.Filename) != ".pdf" {
		return util.ErrorJSON(ctx, errors.New("file extension must be pdf"))
	}

	// Save file in storage folder
	fileName := "usercv" + fmt.Sprint(id)
	ctx.SaveFile(file, "./storage/users/"+fileName+".pdf")
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}
	logger.Logger.Printf("user's cv saved. Name: %s | Size: %d", file.Filename, file.Size)

	// Save fileName in DB
	c.Model.UpdateCV(id, fileName)

	return util.WriteJSON(ctx, http.StatusOK, true, "OK")
}
func (c *Controller) Ban(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		return util.ErrorJSON(ctx, err)
	}

	var req payloads.BanRequest
	err = ctx.BodyParser(&req)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	err = c.Model.Ban(id, req.BanExpire)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, true, "OK")
}
func (c *Controller) Unban(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		return util.ErrorJSON(ctx, err)
	}

	err = c.Model.Unban(id)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, true, "OK")
}
func (c *Controller) Restore(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		logger.Logger.Print(errors.New("invalid id parameter"))
		return util.ErrorJSON(ctx, err)
	}

	err = c.Model.Restore(id)
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
