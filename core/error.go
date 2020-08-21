package core

import "errors"

var (
	//ErrConflict Error when update item that has been modified
	ErrConflict = errors.New("Data has been modified")

	//ErrNotFound Error when item is not found
	ErrNotFound = errors.New("Data was not found")

	//ErrBadRequest Error when data given is not valid on update or insert
	ErrBadRequest = errors.New("Given data is not valid")

	//ErrFailedToCast Error when failed to casting data
	ErrFailedToCast = errors.New("Failed to cast data")

	//ErrZeroAffected Data not found
	ErrZeroAffected = errors.New("No record affected")
)
