package tracer

import (
	"google.golang.org/grpc/codes"
	grpc_trace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

// DdogNonErrorCodeInterceptorOption
// Adjust the list of non-error status codes based on your service's specific needs
var DdogNonErrorCodeInterceptorOption = grpc_trace.NonErrorCodes(
	codes.OK,
	codes.Canceled,
	codes.InvalidArgument,
	codes.NotFound,
	codes.AlreadyExists,
	codes.PermissionDenied,
	codes.ResourceExhausted,
	codes.FailedPrecondition,
	codes.Aborted,
	codes.OutOfRange,
	codes.Unauthenticated,
)
