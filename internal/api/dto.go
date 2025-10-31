package api

type PageEntityDTO interface {
	IsPageEntityDTO() bool // Use in CreatePaginationResult
	GetID() string         // Use in CursorPagination
}
