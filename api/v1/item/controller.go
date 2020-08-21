package item

import (
	"net/http"
	"sample-order/api"
	"sample-order/api/v1/item/request"
	"sample-order/api/v1/item/response"
	"sample-order/core"
	"sample-order/core/item"

	v10 "github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

//Controller Get item API controller
type Controller struct {
	service   item.Service
	validator *v10.Validate
}

//NewController Construct item API controller
func NewController(service item.Service) Controller {
	return Controller{
		service,
		v10.New(),
	}
}

//RegisterPath Register item controller path
func (controller *Controller) RegisterPath(e *echo.Echo) {
	e.GET("/v1/items/:id", controller.getItemByID)
	e.GET("/v1/items/tag/:tag", controller.findItemByTag)

	e.POST("/v1/items", controller.createNewItem)
	e.PUT("/v1/items/:id", controller.updateItem)
}

//Get item by id
func (controller *Controller) getItemByID(c echo.Context) error {
	ID := c.Param("id")
	item, err := controller.service.GetItemByID(ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.NewInternalServerErrorResponse())
	} else if item == nil {
		return c.JSON(http.StatusNotFound, api.NewNotFoundResponse())
	}

	response := response.NewGetItemByIDResponse(item)
	return c.JSON(http.StatusOK, response)
}

//Find items by tag
func (controller *Controller) findItemByTag(c echo.Context) error {
	tag := c.Param("tag")
	items, err := controller.service.GetItemsByTag(tag)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.NewInternalServerErrorResponse())
	}

	response := response.NewGetItemByTagResponse(items)
	return c.JSON(http.StatusOK, response)
}

//Create new item
func (controller *Controller) createNewItem(c echo.Context) error {
	createItemRequest := new(request.CreateItemRequest)

	if err := c.Bind(createItemRequest); err != nil {
		return c.JSON(http.StatusBadRequest, api.NewBadRequestResponse())
	}

	ID, err := controller.service.CreateItem(*createItemRequest.ToUpsertItemSpec(), "creator")

	if err != nil {
		if err == core.ErrBadRequest {
			return c.JSON(http.StatusBadRequest, api.NewBadRequestResponse())
		}
		return c.JSON(http.StatusInternalServerError, api.NewInternalServerErrorResponse())
	}

	response := response.NewCreateNewItemResponse(ID)
	return c.JSON(http.StatusCreated, response)
}

//Update existing item
func (controller *Controller) updateItem(c echo.Context) error {
	updateItemRequest := new(request.UpdateItemRequest)

	if err := c.Bind(updateItemRequest); err != nil {
		return c.JSON(http.StatusBadRequest, api.NewBadRequestResponse())
	}

	err := controller.validator.Struct(updateItemRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.NewBadRequestResponse())
	}

	err = controller.service.UpdateItem(
		c.Param("id"),
		*updateItemRequest.ToUpsertItemSpec(),
		updateItemRequest.Version,
		"updater")

	if err != nil {
		if err == core.ErrNotFound {
			return c.JSON(http.StatusNotFound, api.NewNotFoundResponse())
		}
		if err == core.ErrConflict {
			return c.JSON(http.StatusConflict, api.NewConflictResponse())
		}
	}

	return c.NoContent(http.StatusNoContent)
}
