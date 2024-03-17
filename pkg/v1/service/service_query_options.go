package service

import (
	"fmt"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/options"
)

// RequestWithFilter is an interface that lets us abstract .pb request types that have a filter field.
type RequestWithFilter interface {
	GetFilter() string
}

// getFilterQueryOption returns a []db.QueryOption with a filtering option if the request has a filter.
func getFilterQueryOption[r RequestWithFilter](request r, fieldName string, opts ...options.QueryOption) []options.QueryOption {
	if request.GetFilter() != "" {
		return append(opts, options.WithField(fieldName, request.GetFilter()))
	}
	return opts
}

/* ----------------------------------- */
/*          - Query Options -          */
/* ----------------------------------- */

func signupQueryOptions(username string) []options.QueryOption {
	return []options.QueryOption{options.WithField("username", username)}
}

func loginQueryOptions(username string) []options.QueryOption {
	return []options.QueryOption{options.WithField("username", username)}
}

func getUserQueryOptions(userId int32) []options.QueryOption {
	return []options.QueryOption{options.WithField("id", fmt.Sprint(userId))}
}

func getUsersQueryOptions(in *usersPB.GetUsersRequest, fieldName string) []options.QueryOption {
	return getFilterQueryOption(in, fieldName)
}
