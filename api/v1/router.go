package v1

import (
	"sample-order/api/v1/item"

	"github.com/labstack/echo"
)

//RegisterVIPath Registera V1 API path
func RegisterVIPath(e *echo.Echo, itemController *item.Controller) {
	e.GET("/v1/items/:id", itemController.GetItemByID)
	e.GET("/v1/items/tag/:tag", itemController.FindItemByTag)

	e.POST("/v1/items", itemController.CreateNewItem)
	e.PUT("/v1/items/:id", itemController.UpdateItem)
}
