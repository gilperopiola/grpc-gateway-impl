package tools

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/pbs"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Paginator -            */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type requestsPaginator struct {
	defaultPage     int32
	defaultPageSize int32
}

func NewRequestsPaginator(defaultPage, defaultPageSize int32) *requestsPaginator {
	return &requestsPaginator{defaultPage, defaultPageSize}
}

// Returns the page and pageSize values from a GRPC Request with pagination methods.
// You can set the default values when you create the requestsPaginator.
func (rp *requestsPaginator) PaginatedRequest(req core.PaginatedRequest) (int, int) {
	page := req.GetPage()
	if page == 0 {
		page = rp.defaultPage
	}

	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = rp.defaultPageSize
	}

	return int(page), int(pageSize)
}

// Our paginated endpoints not only take in the page and pageSize from the request but also return
// a *pbs.PaginationInfo with the response, allowing the caller to know the amount of pages that exist for any resource.
func (rp *requestsPaginator) PaginatedResponse(currentPage, pageSize, totalRecords int) *pbs.PaginationInfo {
	totalPages := totalRecords / pageSize
	if totalRecords%pageSize > 0 {
		totalPages++
	}
	return &pbs.PaginationInfo{Current: int32(currentPage), Total: int32(totalPages)}
}
