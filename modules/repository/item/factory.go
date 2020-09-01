package item

import (
	"sample-order/business/item"
	"sample-order/util"
)

//RepositoryFactory Will return business.item.Repository based on active database connection
func RepositoryFactory(dbCon *util.DatabaseConnection) item.Repository {
	var itemRepo item.Repository

	if dbCon.Driver == util.MySQL {
		itemRepo = NewMySQLRepository(dbCon.MySQLDB)
	} else if dbCon.Driver == util.MongoDB {
		itemRepo = NewMongoDBRepository(dbCon.MongoDB)
	}

	return itemRepo
}
