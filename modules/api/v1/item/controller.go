package item

import (
	"net/http"
	"sample-order/business"
	itemBusiness "sample-order/business/item"
	"sample-order/modules/api/common"
	"sample-order/modules/api/v1/item/request"
	"sample-order/modules/api/v1/item/response"

	v10 "github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

//Controller Get item API controller
type Controller struct {
	service   itemBusiness.Service
	validator *v10.Validate
}

//NewController Construct item API controller
func NewController(service itemBusiness.Service) *Controller {
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
		return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
	} else if item == nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	response := response.NewGetItemByIDResponse(*item)
	return c.JSON(http.StatusOK, response)
}

//FindItemByTag Find item by tag echo handler
func (controller *Controller) FindItemByTag(c echo.Context) error {
	tag := c.Param("tag")
	items, err := controller.service.GetItemsByTag(tag)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
	}

	response := response.NewGetItemByTagResponse(items)
	return c.JSON(http.StatusOK, response)
}

//CreateNewItem Create new item echo handler
func (controller *Controller) CreateNewItem(c echo.Context) error {
	createItemRequest := new(request.CreateItemRequest)

	if err := c.Bind(createItemRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	ID, err := controller.service.CreateItem(*createItemRequest.ToUpsertItemSpec(), "creator")

	if err != nil {
		if err == business.ErrInvalidSpec {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}
		return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
	}

	response := response.NewCreateNewItemResponse(ID)
	return c.JSON(http.StatusCreated, response)
}

//UpdateItem update item echo handler
func (controller *Controller) UpdateItem(c echo.Context) error {
	updateItemRequest := new(request.UpdateItemRequest)

	if err := c.Bind(updateItemRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	err := controller.validator.Struct(updateItemRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	err = controller.service.UpdateItem(
		c.Param("id"),
		*updateItemRequest.ToUpsertItemSpec(),
		updateItemRequest.Version,
		"updater")

	if err != nil {
		if err == business.ErrNotFound {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}
		if err == business.ErrHasBeenModified {
			return c.JSON(http.StatusConflict, common.NewConflictResponse())
		}
	}

	return c.NoContent(http.StatusNoContent)
}
