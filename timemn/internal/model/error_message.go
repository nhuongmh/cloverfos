package model

import "errors"

var (
	ErrNotImplemented          = errors.New("method not implemented")
	ErrNoData                  = errors.New("no data")
	ErrNoMoreDataAvailable     = errors.New("no more data available to process")
	ErrServiceIsNotInitialized = errors.New("service is not initialized")
	ErrInvalidInput            = errors.New("invalid input")
	ErrSystemCalculationError  = errors.New("system calculation cerror")
)
