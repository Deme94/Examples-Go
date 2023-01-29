package auth

import (
	"API-REST/api-gateway/controllers/user/auth/payloads"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/conf"
	"API-REST/services/database/postgres/models/user"
	"API-REST/services/logger"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

func (c *Controller) Login(ctx *fiber.Ctx) error {
	var req payloads.LoginRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	var u *user.User
	if req.Email != "" {
		u, err = c.Model.GetByEmailWithPassword(req.Email)
		if err != nil {
			return util.ErrorJSON(ctx, err, http.StatusUnauthorized)
		}
	} else {
		u, err = c.Model.GetByUsernameWithPassword(req.Username)
		if err != nil {
			return util.ErrorJSON(ctx, err, http.StatusUnauthorized)
		}
	}

	if u.DeletedAt != nil {
		return util.ErrorJSON(ctx, errors.New("user deleted"), http.StatusUnauthorized)
	}
	if u.BanDate != nil {
		if u.BanExpire.Before(time.Now()) {
			c.Model.Unban(u.ID)
		} else {
			return util.WriteJSON(ctx, http.StatusUnauthorized, payloads.LoginResponse{BanExpire: u.BanExpire}, "error")
		}
	}

	hashedPassword := u.Password

	err = c.compareHashAndPassword(hashedPassword, req.Password)
	if err != nil {
		return util.ErrorJSON(ctx, err, http.StatusUnauthorized)
	}

	// Generate jwt token after successful login
	token, err := c.generateJwtToken(fmt.Sprint(u.ID), conf.Env.GetString("JWT_SECRET"))
	if err != nil {
		return util.ErrorJSON(ctx, err, http.StatusNotImplemented)
	}

	return util.WriteJSON(ctx, http.StatusOK, payloads.LoginResponse{ID: u.ID, Token: string(token)}, "response")
}

func (c *Controller) Get(ctx *fiber.Ctx) error {
	claimerID := ctx.Locals("Claimer-ID").(int)

	u, err := c.Model.Get(claimerID)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	userResponse := payloads.GetResponse{
		ID:                 u.ID,
		CreatedAt:          u.CreatedAt,
		Username:           u.Username,
		Email:              u.Email,
		Nick:               u.Nick,
		FirstName:          u.FirstName,
		LastName:           u.LastName,
		Phone:              u.Phone,
		Address:            u.Address,
		LastPasswordChange: u.LastPasswordChange,
		VerifiedMail:       u.VerifiedMail,
		VerifiedPhone:      u.VerifiedPhone,
	}

	return util.WriteJSON(ctx, http.StatusOK, userResponse, "user")
}
func (c *Controller) GetPhoto(ctx *fiber.Ctx) error {
	claimerID := ctx.Locals("Claimer-ID").(int)

	imageName, err := c.Model.GetPhoto(claimerID)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	// Get webp file
	f, err := os.OpenFile("./storage/users/"+imageName+".webp", os.O_RDWR, 0644)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}
	defer f.Close()
	// Decode webp file to image
	image, err := webp.Decode(f, &decoder.Options{})
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}
	// Encode image into buffer
	var buf bytes.Buffer
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	if err != nil {
		log.Fatalln(err)
	}
	webp.Encode(&buf, image, options)
	// Get bytes and encode to base64
	imageBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return util.WriteJSON(ctx, http.StatusOK, imageBase64, "photo")
}
func (c *Controller) GetCV(ctx *fiber.Ctx) error {
	claimerID := ctx.Locals("Claimer-ID").(int)

	cvName, err := c.Model.GetCV(claimerID)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	filePath := "./storage/users/" + cvName + ".pdf"
	return ctx.SendFile(filePath)
}
func (c *Controller) Update(ctx *fiber.Ctx) error {
	claimerID := ctx.Locals("Claimer-ID").(int)

	var req payloads.UpdateRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	var u user.User
	u.ID = claimerID
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
func (c *Controller) ChangePassword(ctx *fiber.Ctx) error {
	claimerID := ctx.Locals("Claimer-ID").(int)

	var req payloads.ChangePasswordRequest

	err := ctx.BodyParser(&req)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	hashedPassword, err := c.Model.GetPassword(claimerID)
	if err != nil {
		return util.ErrorJSON(ctx, err, http.StatusInternalServerError)
	}

	err = c.compareHashAndPassword(hashedPassword, req.OldPassword)
	if err != nil {
		return util.ErrorJSON(ctx, err, http.StatusUnauthorized)
	}

	hashedNewPassword, err := c.hashPassword(req.NewPassword)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	err = c.Model.UpdatePassword(claimerID, hashedNewPassword)
	if err != nil {
		return util.ErrorJSON(ctx, err, http.StatusInternalServerError)
	}

	return util.WriteJSON(ctx, http.StatusOK, true, "OK")
}
func (c *Controller) ResetPassword(ctx *fiber.Ctx) error {
	var req payloads.ResetPasswordRequest

	err := ctx.BodyParser(&req)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	u, err := c.Model.GetByEmailWithPassword(req.Email)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	if u.DeletedAt != nil {
		return util.ErrorJSON(ctx, errors.New("user deleted"), http.StatusUnauthorized)
	}

	password := c.generateRandomPassword()

	err = c.Model.UpdatePassword(u.ID, password)
	if err != nil {
		return util.ErrorJSON(ctx, err, http.StatusInternalServerError)
	}

	// Send password to email
	fmt.Println(password)
	// ...

	return util.WriteJSON(ctx, http.StatusOK, true, "OK")
}
func (c *Controller) UpdateRoles(ctx *fiber.Ctx) error {
	claimerID := ctx.Locals("Claimer-ID").(int)

	var req payloads.UpdateRolesRequest

	err := ctx.BodyParser(&req)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	err = c.Model.UpdateRoles(claimerID, req.RoleIDs...)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, true, "OK")
}
func (c *Controller) UpdatePhoto(ctx *fiber.Ctx) error {
	claimerID := ctx.Locals("Claimer-ID").(int)

	var req payloads.UpdatePhotoRequest

	err := ctx.BodyParser(&req)
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
	imageName := "user" + fmt.Sprint(claimerID)
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

	err = c.Model.UpdatePhoto(claimerID, imageName)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, true, "OK")
}
func (c *Controller) UpdateCV(ctx *fiber.Ctx) error {
	claimerID := ctx.Locals("Claimer-ID").(int)

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
	fileName := "usercv" + fmt.Sprint(claimerID)
	ctx.SaveFile(file, "./storage/users/"+fileName+".pdf")
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}
	logger.Logger.Printf("user's cv saved. Name: %s | Size: %d", file.Filename, file.Size)

	// Save fileName in DB
	c.Model.UpdateCV(claimerID, fileName)

	return util.WriteJSON(ctx, http.StatusOK, true, "OK")
}
func (c *Controller) Delete(ctx *fiber.Ctx) error {
	claimerID := ctx.Locals("Claimer-ID").(int)

	err := c.Model.Delete(claimerID)
	if err != nil {
		return util.ErrorJSON(ctx, err)
	}

	return util.WriteJSON(ctx, http.StatusOK, true, "OK")
}
