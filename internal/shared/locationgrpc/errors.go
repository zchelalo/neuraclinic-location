package locationgrpc

import (
	"errors"

	locationerrors "github.com/zchelalo/neuraclinic-location/internal/shared/locationerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapError(err error) error {
	switch {
	case errors.Is(err, locationerrors.ErrNotFound):
		return status.Error(codes.NotFound, "not found")
	case errors.Is(err, locationerrors.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, "invalid input")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
