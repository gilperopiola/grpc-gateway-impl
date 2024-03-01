package interceptors

import (
	"errors"
	"testing"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/bufbuild/protovalidate-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* ----------------------------------- */
/*        - Interceptor Tests -        */
/* ----------------------------------- */

func TestFromValidationErrToGRPCInvalidArgErr(t *testing.T) {
	tests := []struct {
		name             string
		validationErr    error
		invalidArgErrMsg string
	}{
		{
			name: "validation_error",
			validationErr: &protovalidate.ValidationError{
				Violations: []*validate.Violation{{FieldPath: "field", Message: "is required"}},
			},
			invalidArgErrMsg: "validation error: field is required",
		},
		{
			name:             "unexpected_error",
			validationErr:    errors.New("unexpected error"),
			invalidArgErrMsg: "unexpected validation error: unexpected error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fromValidationErrToGRPCInvalidArgErr(tt.validationErr)
			grpcStatus, _ := status.FromError(err)

			assert.Equal(t, codes.InvalidArgument, grpcStatus.Code())
			assert.Equal(t, tt.invalidArgErrMsg, grpcStatus.Message())
		})
	}
}

func TestNewProtoValidator(t *testing.T) {
	validator := NewProtoValidator()
	assert.NotNil(t, validator)
}
