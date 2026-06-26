package locationgrpc

import (
	"context"
	"testing"

	"github.com/zchelalo/neuraclinic-location/internal/shared/appctx"
	"github.com/zchelalo/neuraclinic-location/internal/shared/i18n"
	locationerrors "github.com/zchelalo/neuraclinic-location/internal/shared/locationerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapErrorLocalizesMessage(t *testing.T) {
	t.Parallel()

	ctx := appctx.WithLanguage(context.Background(), i18n.Spanish)
	err := MapError(ctx, locationerrors.ErrInvalidInput)
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected grpc status, got %v", err)
	}
	if st.Code() != codes.InvalidArgument {
		t.Fatalf("status code = %s, want %s", st.Code(), codes.InvalidArgument)
	}
	if st.Message() != "entrada invalida" {
		t.Fatalf("status message = %q", st.Message())
	}
}
