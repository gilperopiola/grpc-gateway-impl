package service

import (
	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/common"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* ----------------------------------- */
/*           - v1 Service -            */
/* ----------------------------------- */

// Service holds every gRPC method of the usersPB.UsersServiceServer.
// It handles the business logic of the API.
type Service interface {
	usersPB.UsersServiceServer
}

// service is our concrete implementation of the Service interface.
// It has an embedded Repository to interact with the database.
type service struct {
	Repo           repository.Repository
	TokenGenerator common.TokenGenerator
	PwdHasher      common.PwdHasher

	*usersPB.UnimplementedUsersServiceServer
}

// NewService returns a new instance of the service.
func NewService(repo repository.Repository, tokenGenerator common.TokenGenerator, pwdHasher common.PwdHasher) *service {
	return &service{
		Repo:           repo,
		TokenGenerator: tokenGenerator,
		PwdHasher:      pwdHasher,
	}
}

/* ----------------------------------- */
/*           - Pagination -            */
/* ----------------------------------- */

// RequestWithPagination is an interface that lets us abstract .pb request types that have pagination methods.
type RequestWithPagination interface {
	GetPage() int32
	GetPageSize() int32
}

// getPaginationValues returns the page and pageSize values from a gRPC Request with pagination methods.
// It defaults to page 1 and pageSize 10 if the values are not set or are invalid.
func getPaginationValues[r RequestWithPagination](req r) (int, int) {
	page, pageSize := req.GetPage(), req.GetPageSize()
	defaultPage, defaultPageSize := 1, 10

	if page == 0 {
		page = int32(defaultPage)
	}

	if pageSize == 0 {
		pageSize = int32(defaultPageSize)
	}

	return int(page), int(pageSize)
}

// makeResponsePagination returns a PaginationInfo object with the current and total pages.
func makeResponsePagination(page, pageSize, totalMatchingRecords int) *usersPB.PaginationInfo {
	totalPages := totalMatchingRecords / pageSize

	if totalMatchingRecords%pageSize > 0 {
		totalPages++
	}
	return &usersPB.PaginationInfo{Current: int32(page), Total: int32(totalPages)}
}

/* ----------------------------------- */
/*         - Service Errors -          */
/* ----------------------------------- */

var grpcUnknownErr = func(msg string, err error) error {
	return status.Errorf(codes.Unknown, "%s: %v", msg, err)
}

var grpcNotFoundErr = func(entity string) error {
	return status.Errorf(codes.NotFound, "%s not found.", entity)
}

var grpcAlreadyExistsErr = func(entity string) error {
	return status.Errorf(codes.AlreadyExists, "%s already exists.", entity)
}

var grpcUnauthenticatedErr = func(reason string) error {
	return status.Errorf(codes.Unauthenticated, reason)
}

var grpcInvalidArgumentErr = func(entity string) error {
	return status.Errorf(codes.InvalidArgument, "invalid %s.", entity)
}
