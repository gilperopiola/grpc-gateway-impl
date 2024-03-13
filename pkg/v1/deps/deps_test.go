package deps

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/bufbuild/protovalidate-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewProtoValidator(t *testing.T) {
	validator := NewValidator()
	assert.NotNil(t, validator)
}

func TestHandleValidationErr(t *testing.T) {
	tests := []struct {
		name          string
		validationErr error
		outErrMsg     string
	}{
		{
			name: "validation_error",
			validationErr: &protovalidate.ValidationError{
				Violations: []*validate.Violation{{FieldPath: "field", Message: "is required"}},
			},
			outErrMsg: fmt.Sprintf(errs.ErrMsgInValidation, "field is required"),
		},
		{
			name:          "unexpected_error",
			validationErr: errors.New("unexpected error"),
			outErrMsg:     fmt.Sprintf(errs.ErrMsgInValidationUnexpected, "unexpected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleValidationErr(tt.validationErr)
			grpcStatus, _ := status.FromError(err)

			assert.Equal(t, codes.InvalidArgument, grpcStatus.Code())
			assert.Equal(t, tt.outErrMsg, grpcStatus.Message())
		})
	}
}
