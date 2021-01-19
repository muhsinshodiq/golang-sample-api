package http

import (
	"sample-order/api/v1/item"

	"github.com/labstack/echo"
)

//RegisterPath Registera V1 API path
func RegisterPath(e *echo.Echo, itemController *item.Controller) {
	if itemController == nil {
		panic("item controller cannot be nil")
	}

	//item
	itemV1 := e.Group("v1/items")
	itemV1.GET("/:id", itemController.GetItemByID)
	itemV1.GET("/tag/:tag", itemController.FindItemByTag)
	itemV1.POST("", itemController.CreateNewItem)
	itemV1.PUT("/:id", itemController.UpdateItem)

	//health check
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(200)
	})
}
