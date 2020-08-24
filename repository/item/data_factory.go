package item

import (
	itemCore "sample-order/core/item"
	"sample-order/libs"
)

//DataRepositoryFactory Will return *itemCore.DataRepository based on active database connection
func DataRepositoryFactory(dbCon *libs.DatabaseConnection) itemCore.DataRepository {
	var itemDataRepo itemCore.DataRepository

	if dbCon.Driver == libs.MySQL {
		itemDataRepo = NewMySQLRepository(dbCon.MySQLDB)
	} else if dbCon.Driver == libs.MongoDB {
		itemDataRepo = NewMongoDBRepository(dbCon.MongoDB)
	}

	return itemDataRepo
}
