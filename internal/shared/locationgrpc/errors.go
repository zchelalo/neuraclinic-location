package locationgrpc

import (
	"context"
	"errors"

	"github.com/zchelalo/neuraclinic-location/internal/shared/appctx"
	"github.com/zchelalo/neuraclinic-location/internal/shared/i18n"
	locationerrors "github.com/zchelalo/neuraclinic-location/internal/shared/locationerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapError(ctx context.Context, err error) error {
	language := appctx.Language(ctx)
	switch {
	case errors.Is(err, locationerrors.ErrNotFound):
		return status.Error(codes.NotFound, i18n.Message(language, i18n.KeyNotFound))
	case errors.Is(err, locationerrors.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, i18n.Message(language, i18n.KeyInvalidInput))
	default:
		return status.Error(codes.Internal, i18n.Message(language, i18n.KeyInternalServerError))
	}
}
