package tools

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var _ core.RequestsValidator = (*protoRequestsValidator)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - Requests Validator -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Validates GRPC requests based on rules written on .proto files.
// It uses the bufbuild/protovalidate library.
type protoRequestsValidator struct {
	*protovalidate.Validator
}

// New instance of *protoValidator. This panics on failure.
func NewRequestsValidator() core.RequestsValidator {
	validator, err := protovalidate.New()
	core.LogFatalIfErr(err, errs.FailedToCreateProtoVal)
	return &protoRequestsValidator{validator}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Wraps the proto validation logic with a GRPC Interceptor.
func (prv protoRequestsValidator) ValidateGRPC(ctx god.Ctx, req any, _ *god.GRPCInfo, handler god.GRPCHandler) (any, error) {
	if err := prv.Validate(req.(protoreflect.ProtoMessage)); err != nil {
		return nil, validationErrToGRPC(err)
	}
	return handler(ctx, req)
}

// Takes a *protovalidate.ValidationError and returns an InvalidArgument(3) GRPC error with its corresponding message.
// Validation errors are always returned as InvalidArgument.
// They get translated to 400 Bad Request on the HTTP error handler.
func validationErrToGRPC(err error) error {
	var validationErr *protovalidate.ValidationError
	if ok := errors.As(err, &validationErr); ok {
		humanFacingMsg := brokenRulesWrap{validationErr.Violations}.HumanFacing(",")
		return status.Error(codes.InvalidArgument, fmt.Sprintf(errs.ValidatingRequest, humanFacingMsg))
	}

	var runtimeErr *protovalidate.RuntimeError
	if ok := errors.As(err, &runtimeErr); ok {
		core.LogUnexpected(runtimeErr)
		return status.Error(codes.InvalidArgument, fmt.Sprintf(errs.ValidatingRequestRuntime, runtimeErr))
	}

	core.LogUnexpected(err)
	return status.Error(codes.InvalidArgument, fmt.Sprintf(errs.ValidatingRequestUnexpected, err))
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// -> To avoid using the word 'violation' in the codebase we wrap the protovalidate.ValidationError and handle
// everything there as a 'broken rule'.

type brokenRulesWrap struct {
	Rules []*validate.Violation
}

type brokenRule validate.Violation

// Returns a string detailing the rules broken.
// This string is the human-facing format in which validation errors translate.
func (wrap brokenRulesWrap) HumanFacing(delimiter string) string {
	out := ""
	for i, rule := range wrap.Rules {
		out += (*brokenRule)(rule).HumanFacing()
		if i < len(wrap.Rules)-1 {
			out += (delimiter + " ")
		}
	}
	return out
}

// Pointer receiver because brokenRule has a mutex.
func (rule *brokenRule) HumanFacing() string {
	if strings.Contains(rule.Message, "match regex pattern") {
		return fmt.Sprintf("%s value has an invalid format", rule.FieldPath) // don't show regex pattern.
	}
	return fmt.Sprintf("%s %s", rule.FieldPath, rule.Message) // simple message.
}
