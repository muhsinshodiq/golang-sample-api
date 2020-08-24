package item

import (
	"net/http"
	"sample-order/api"
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
func NewController(service item.Service) *Controller {
	return &Controller{
		service,
		v10.New(),
	}
}

//GetItemByID Get item by ID echo handler
func (controller *Controller) GetItemByID(c echo.Context) error {
	ID := c.Param("id")
	item, err := controller.service.GetItemByID(ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.NewInternalServerErrorResponse())
	} else if item == nil {
		return c.JSON(http.StatusNotFound, api.NewNotFoundResponse())
	}

	response := NewGetItemByIDResponse(item)
	return c.JSON(http.StatusOK, response)
}

//FindItemByTag Find item by tag echo handler
func (controller *Controller) FindItemByTag(c echo.Context) error {
	tag := c.Param("tag")
	items, err := controller.service.GetItemsByTag(tag)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.NewInternalServerErrorResponse())
	}

	response := NewGetItemByTagResponse(items)
	return c.JSON(http.StatusOK, response)
}

//CreateNewItem Create new item echo handler
func (controller *Controller) CreateNewItem(c echo.Context) error {
	createItemRequest := new(CreateItemRequest)

	if err := c.Bind(createItemRequest); err != nil {
		return c.JSON(http.StatusBadRequest, api.NewBadRequestResponse())
	}

	ID, err := controller.service.CreateItem(*createItemRequest.ToUpsertItemSpec(), "creator")

	if err != nil {
		if err == item.ErrInvalidSpec {
			return c.JSON(http.StatusBadRequest, api.NewBadRequestResponse())
		}
		return c.JSON(http.StatusInternalServerError, api.NewInternalServerErrorResponse())
	}

	response := NewCreateNewItemResponse(ID)
	return c.JSON(http.StatusCreated, response)
}

//UpdateItem update item echo handler
func (controller *Controller) UpdateItem(c echo.Context) error {
	updateItemRequest := new(UpdateItemRequest)

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
		if err == item.ErrNotFound {
			return c.JSON(http.StatusNotFound, api.NewNotFoundResponse())
		}
		if err == item.ErrDataHasBeenModified {
			return c.JSON(http.StatusConflict, api.NewConflictResponse())
		}
	}

	return c.NoContent(http.StatusNoContent)
}
