package model

type StorageResponse struct {
	Key string
	Value int64
	ActionTimestamp string
}

type ErrorResponse struct {
	ErrorMessage string
	TransactionTs string
}
