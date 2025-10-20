package model

type Sort struct {
	Column    string
	Direction string
}

type PaginationResponse struct {
	Limit      int32
	Page       int32
	TotalRows  int64
	TotalPages int32
}
