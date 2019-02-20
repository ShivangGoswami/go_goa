package main

import (
	"database/sql"
	"ryan/app"
	"ryan/respositories"
	"ryan/util/crypto"
	"ryan/util/jwt"

	"github.com/goadesign/goa"
)

// AuthenticationController implements the authentication resource.
type AuthenticationController struct {
	*goa.Controller
	*sql.DB
}

// NewAuthenticationController creates a authentication controller.
func NewAuthenticationController(service *goa.Service, db *sql.DB) *AuthenticationController {
	return &AuthenticationController{
		Controller: service.NewController("AuthenticationController"),
		DB:         db,
	}
}

// Login runs the login action.
func (c *AuthenticationController) Login(ctx *app.LoginAuthenticationContext) error {
	// AuthenticationController_Login: start_implement
	payload := ctx.Payload
	u, err := respositories.GetUserByEmail(c.DB, payload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.BadRequest(goa.ErrBadRequest("Invalid email or password"))
		}
		c.Service.LogError("Login User", "err", err)
		return ctx.InternalServerError()
	}
	hashedPassword := crypto.HashPassword(payload.Password, *u.Salt)
	if *u.Password != hashedPassword {
		return ctx.BadRequest(goa.ErrBadRequest("Invalid email or password"))
	}
	token, err := jwt.CreateJWTToken(*u.Email)
	if err != nil {
		c.Service.LogError("Login User", "err", err)
		return ctx.InternalServerError()
	}
	// Put your logic here

	res := &app.Token{
		Token: &token,
	}
	return ctx.OK(res)
	// AuthenticationController_Login: end_implement
}

// Register runs the register action.
func (c *AuthenticationController) Register(ctx *app.RegisterAuthenticationContext) error {
	// AuthenticationController_Register: start_implement
	payload := ctx.Payload
	exists, err := respositories.CheckEmailExists(c.DB, payload.Email)
	if err != nil {
		c.Service.LogError("Register User", "err", err)
		return ctx.InternalServerError()
	}
	if exists {
		return ctx.BadRequest(goa.ErrBadRequest("Email already exists"))
	}
	err = respositories.AddUserToDatabase(c.DB, payload.FirstName, payload.LastName, payload.Email, payload.Password)
	if err != nil {
		c.Service.LogError("Register User", "err", err)
		return ctx.InternalServerError()
	}
	token, err := jwt.CreateJWTToken(payload.Email)
	if err != nil {
		c.Service.LogError("Register User", "err", err)
		return ctx.InternalServerError()
	}
	res := &app.Token{Token: &token}
	return ctx.OK(res)
	// Put your logic here
	// AuthenticationController_Register: end_implement
}
