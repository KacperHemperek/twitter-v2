package services

import "errors"

var (
	ErrResourceNotFound     = errors.New("resource not found")
	ErrInvalidQueryResponse = errors.New("invalid query response")
)
