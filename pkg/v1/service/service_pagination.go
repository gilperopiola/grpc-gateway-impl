package service

import usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"

/* ----------------------------------- */
/*           - Pagination -            */
/* ----------------------------------- */

// RequestWithPagination is an interface that lets us abstract .pb request types that have pagination methods.
type RequestWithPagination interface {
	GetPage() int32
	GetPageSize() int32
}

// getPaginationValues returns the page and pageSize values from a gRPC Request with pagination methods.
// It defaults to page 1 and pageSize 10 if the values are not set.
func getPaginationValues[r RequestWithPagination](req r) (int, int) {
	defaultPage := 1      // T0D0 -> Config var.
	defaultPageSize := 10 // T0D0 -> Config var.

	page := req.GetPage()
	if page == 0 {
		page = int32(defaultPage)
	}

	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = int32(defaultPageSize)
	}

	return int(page), int(pageSize)
}

// makeResponsePagination returns a *usersPB.PaginationInfo with the current and total pages.
func makeResponsePagination(page, pageSize, totalMatchingRecords int) *usersPB.PaginationInfo {
	totalPages := totalMatchingRecords / pageSize

	if totalMatchingRecords%pageSize > 0 {
		totalPages++
	}

	return &usersPB.PaginationInfo{
		Current: int32(page),
		Total:   int32(totalPages),
	}
}
