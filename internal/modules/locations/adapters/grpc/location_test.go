package grpc

import (
	"errors"
	"testing"

	locationv1 "github.com/zchelalo/neuraclinic-location/gen/go/location/v1"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/domain"
	locationerrors "github.com/zchelalo/neuraclinic-location/internal/shared/locationerrors"
)

func TestAdminAreaTypeFilterFromProtoUsesEnum(t *testing.T) {
	adminAreaType := locationv1.AdminAreaType_ADMIN_AREA_TYPE_MUNICIPALITY

	got, err := adminAreaTypeFilterFromProto(&locationv1.ListAdminAreasRequest{
		AdminAreaType: &adminAreaType,
	})
	if err != nil {
		t.Fatalf("adminAreaTypeFilterFromProto() error = %v", err)
	}
	if got != domain.AdminAreaTypeMunicipality {
		t.Fatalf("adminAreaTypeFilterFromProto() = %q, want %q", got, domain.AdminAreaTypeMunicipality)
	}
}

func TestAdminAreaTypeFilterFromProtoWithoutEnumReturnsEmptyFilter(t *testing.T) {
	got, err := adminAreaTypeFilterFromProto(&locationv1.ListAdminAreasRequest{})
	if err != nil {
		t.Fatalf("adminAreaTypeFilterFromProto() error = %v", err)
	}
	if got != "" {
		t.Fatalf("adminAreaTypeFilterFromProto() = %q, want empty filter", got)
	}
}

func TestAdminAreaTypeFromProtoRejectsUnknownValue(t *testing.T) {
	_, err := adminAreaTypeFromProto(locationv1.AdminAreaType(99))
	if !errors.Is(err, locationerrors.ErrInvalidInput) {
		t.Fatalf("adminAreaTypeFromProto() error = %v, want ErrInvalidInput", err)
	}
}

func TestAdminAreaTypeToProto(t *testing.T) {
	if got := adminAreaTypeToProto(domain.AdminAreaTypeState); got != locationv1.AdminAreaType_ADMIN_AREA_TYPE_STATE {
		t.Fatalf("adminAreaTypeToProto(state) = %s", got)
	}
	if got := adminAreaTypeToProto(domain.AdminAreaTypeMunicipality); got != locationv1.AdminAreaType_ADMIN_AREA_TYPE_MUNICIPALITY {
		t.Fatalf("adminAreaTypeToProto(municipality) = %s", got)
	}
	if got := adminAreaTypeToProto("province"); got != locationv1.AdminAreaType_ADMIN_AREA_TYPE_UNSPECIFIED {
		t.Fatalf("adminAreaTypeToProto(unknown) = %s", got)
	}
}
