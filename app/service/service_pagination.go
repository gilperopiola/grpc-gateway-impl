package service

import "github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - Pagination -            */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// PaginatedRequest is an interface that lets us abstract .pb request types that have pagination methods.
type PaginatedRequest interface {
	GetPage() int32
	GetPageSize() int32
}

// getPaginationFromRequest returns the page and pageSize values from a GRPC Request with pagination methods.
// It defaults to page 1 and pageSize 10 if the values are not set.
func getPaginationFromRequest[r PaginatedRequest](req r) (int, int) {
	defaultPage := int32(1)      // T0D0 -> Config var.
	defaultPageSize := int32(10) // T0D0 -> Config var.

	page := req.GetPage()
	if page == 0 {
		page = defaultPage
	}

	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = defaultPageSize
	}

	return int(page), int(pageSize)
}

// Part of the Response of a GetMany endpoint (like GetUsers).
// Has Current and Total Pages.
func newResponsePagination(currentPage, pageSize, totalRecords int) *pbs.PaginationInfo {
	totalPages := totalRecords / pageSize
	if totalRecords%pageSize > 0 {
		totalPages++
	}

	return &pbs.PaginationInfo{
		Current: int32(currentPage),
		Total:   int32(totalPages),
	}
}
