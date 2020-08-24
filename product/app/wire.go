package main

import (
	itemControllerV1 "sample-order/api/v1/item"
	itemCore "sample-order/core/item"
	"sample-order/libs"
	itemRepo "sample-order/repository/item"

	"github.com/google/wire"
)

//InitializeItemControllerV1 Initialize item service
func InitializeItemControllerV1(dbCon *libs.DatabaseConnection) *itemControllerV1.Controller {
	wire.Build(
		itemControllerV1.NewController,
		itemCore.NewService,
		itemRepo.DataRepositoryFactory)

	return nil
}
