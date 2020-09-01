package business

import "errors"

var (
	//ErrHasBeenModified Error when update item that has been modified
	ErrHasBeenModified = errors.New("Data has been modified")

	//ErrNotFound Error when item is not found
	ErrNotFound = errors.New("Data was not found")

	//ErrInvalidSpec Error when data given is not valid on update or insert
	ErrInvalidSpec = errors.New("Given spec is not valid")

	//ErrZeroAffected Data not found
	ErrZeroAffected = errors.New("No record affected")
)
