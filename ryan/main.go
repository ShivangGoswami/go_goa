//go:generate goagen bootstrap -d ryan/design

package main

import (
	"ryan/app"
	"ryan/util/database"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
)

func main() {
	// Create service
	service := goa.New("Authentication API")
	db, err := database.Connect()
	if err != nil {
		service.LogError("startup", "err", err)
	}
	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "authentication" controller
	c := NewAuthenticationController(service, db)
	app.MountAuthenticationController(service, c)

	// Start service
	if err := service.ListenAndServe(":8000"); err != nil {
		service.LogError("startup", "err", err)
	}

}
