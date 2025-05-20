package api

import (
	"bridgebot/internal/api/bridge_swap"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewServer(swapServer *bridge_swap.SwapServer) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	e.Validator = &CustomValidator{validator: validator.New()}
	// todo: Use ECHO Groupt 
	e.POST("/v1/getQuote",  swapServer.HandleQuote)
	e.POST("/v1/swap", swapServer.HandleSwap)

	return e
}
